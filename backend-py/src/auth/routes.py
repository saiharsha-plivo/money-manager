from fastapi import APIRouter, BackgroundTasks, Depends, status, HTTPException, Response, Request
from fastapi.responses import JSONResponse
from .schemas import UserSignUp, ForgotPassword
from src.db.main import get_session
from sqlmodel.ext.asyncio.session import AsyncSession
from .service import AuthService
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from .utils import verify_password, create_jwt_token, decode_jwt_token, generate_url_safe_token, decode_url_safe_token
from src.db.models import User
from src.config import Config
from datetime import timedelta
from src.logger import get_logger
from src.mail import EmailSchema, schedule_email
from .utils import generate_passwd_hash


version = "v1"
version_prefix =f"/api/{version}"

logger = get_logger("auth-router")
auth_router = APIRouter()
auth_service = AuthService()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl=f"{version_prefix}/auth/login")

async def get_current_user(token: str = Depends(oauth2_scheme),session = Depends(get_session)):
  logger.info("request for get current user")
  username = decode_jwt_token(token)
  if username is None:
    raise HTTPException(status_code=401, detail="Invalid token")

  user = await auth_service.check_user_exits(username, session)
  if user is None:
    raise HTTPException(status_code=401, detail="User not found")
  
  return user

@auth_router.post("/sign-up", status_code=status.HTTP_201_CREATED)
async def sign_up_user(user_signup_data: UserSignUp,background_tasks: BackgroundTasks,session: AsyncSession = Depends(get_session)):
  username = user_signup_data.username
  logger.debug(f"request for signup username : {username}")
  user = await auth_service.check_user_exits(username, session)

  if user is not None:
    raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Username already exists")

  newuser = await auth_service.create_user(user_signup_data, session)

  verify_token = generate_url_safe_token(user_signup_data.email)

  link = f"http://{Config.DOMAIN}/api/v1/auth/verify/{verify_token}"

  mail_message = EmailSchema(
    subject="Money Manager Verification",
    email=[user_signup_data.email],
    message=f"Please click the link to verify the user {user_signup_data.username} : {link}",
  )

  schedule_email(background_tasks, mail_message)

  return {
    "user": newuser,
    "detail": "Account Created! Check email to verify your account",
  }
  

@auth_router.get("/verify/{verify_token}")
async def verify_token(verify_token: str,session: AsyncSession = Depends(get_session)):
  email = decode_url_safe_token(verify_token)
  
  if email is None:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail="Invalid token")
  
  user = await auth_service.get_user_by_email(email,session)
  
  if user is None:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail=f"User with email {email} not found , please sign up again")
  
  if user.verified:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail="User already verified")
  
  await auth_service.update_user(user,{"verified":True},session)
  
  return {"detail": "User verified successfully"}
    
    
@auth_router.post("/login",status_code=status.HTTP_200_OK)
async def login(response: Response,form_data: OAuth2PasswordRequestForm = Depends(),session: AsyncSession = Depends(get_session)):
  user : User = await auth_service.check_user_exits(form_data.username,session)
  
  if user is None :
    raise HTTPException(status_code=status.HTTP_404_NOT_FOUND,detail="user does not exists")
  
  logger.debug(f"request for signup username : {user.username}")
  password_check = verify_password(password=form_data.password,hash=user.password_hash)
  
  if not password_check : 
    raise HTTPException(status_code=status.HTTP_403_FORBIDDEN,detail="password is incorrect")
  
  access_token = create_jwt_token({"sub": user.username},timedelta(minutes=Config.ACCESS_TOKEN_EXPIRE_MINUTES))
  refresh_token = create_jwt_token({"sub": user.username},timedelta(days=Config.REFRESH_TOKEN_EXPIRE_DAYS))
  
  response.set_cookie(
    key="refresh_token",
    value=refresh_token,
    httponly=True,
    secure=False,
    samesite="lax",
    max_age=Config.REFRESH_TOKEN_EXPIRE_DAYS * 86400,
    path="/"
  )
  
  return {"access_token": access_token, "token_type": "bearer"}

@auth_router.get("/logout", status_code=status.HTTP_200_OK)
async def logout(request: Request, response: Response):
  logger.debug("request for logout username")
  cookies = request.cookies
  refresh_token = cookies.get("refresh_token")

  if not refresh_token:
    return JSONResponse(
        status_code=status.HTTP_204_NO_CONTENT,
        content={"detail": "No refresh token cookie found"}
    )

  response.delete_cookie(key="refresh_token", path="/")
  return {"detail": "User logged out successfully"}

@auth_router.post("/refresh")
def refresh(request: Request):
  logger.debug("request for refresh token")
  refresh_token = request.cookies.get("refresh_token")
  username = decode_jwt_token(refresh_token)
  
  if username is None :
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED,detail="Refresh Token is invalid")
  
  access_token = create_jwt_token({"sub": username}, timedelta(minutes=Config.ACCESS_TOKEN_EXPIRE_MINUTES))
  
  return {"access_token": access_token, "token_type": "bearer"}

  
@auth_router.get("/users/me")
async def current_user(user=Depends(get_current_user)):
  logger.debug("request to get current user")
  return user


@auth_router.post("/forgot-password")
async def forgot_password(username: str,background_tasks: BackgroundTasks,session: AsyncSession = Depends(get_session)):
  user = await auth_service.check_user_exits(username,session)
  
  logger.debug(f"request for forgot password username : {username}")
  
  if user is None:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail=f"User with username {username} not found")
  
  verify_token = generate_url_safe_token(user.email)
  
  link = f"http://{Config.DOMAIN}/api/v1/auth/reset-password/{verify_token}"
  
  mail_message = EmailSchema(
    subject="Money Manager Password Reset",
    email=[user.email],
    message=f"Please click the link to reset the password for {user.username} : {link}",
  )
  
  schedule_email(background_tasks, mail_message)
  logger.debug(f"password reset link sent to email : {user.email}")
  
  return {"detail": "Password reset link sent to email"}
  
  
@auth_router.post("/reset-password/{reset_token}")
async def reset_password(reset_token: str,password_data: ForgotPassword,session: AsyncSession = Depends(get_session)):
  email = decode_url_safe_token(reset_token)
  logger.debug(f"request for reset password email : {email}")
  if email is None: 
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail="Invalid token")
  
  user = await auth_service.get_user_by_email(email,session)
  
  if user is None:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail=f"User with email {email} not found")
  
  if password_data.password != password_data.confirm_password:
    raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,detail="Passwords do not match")
  
  password_hash = generate_passwd_hash(password_data.password)
  await auth_service.update_user(user,{"password_hash":password_hash},session)
  logger.debug(f"password reset successfully for user : {user.username}")
  
  return {"detail": "Password reset successfully"}
  
  
  
  
  
  