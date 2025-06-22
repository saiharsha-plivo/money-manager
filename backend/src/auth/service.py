from sqlmodel.ext.asyncio.session import AsyncSession
from src.db.models import User
from sqlmodel import select
from .schemas import UserSignUp
from .utils import generate_passwd_hash

class AuthService:
  async def check_user_exits(self, username: str, session: AsyncSession) -> User :
    statement = (
      select(User).where(
        User.username == username
      )
    )
    result = await session.exec(statement)
    user = result.first()
    return user
  
  async def get_user_by_email(self, email: str, session: AsyncSession) -> User :
    statement = (
      select(User).where(
        User.email == email
      )
    )
    result = await session.exec(statement)
    user = result.first()
    return user
  
  async def create_user(self,user_details: UserSignUp, session: AsyncSession):
    password = user_details.password
    password_hash = generate_passwd_hash(password)
    
    newuser = User(
      username=user_details.username,
      role="user",
      password_hash=password_hash,
      email= user_details.email
    )
    
    session.add(newuser)
    await session.commit()
    
    return newuser
  
  async def update_user(self, user:User , user_data: dict,session:AsyncSession):
    for k, v in user_data.items():
      setattr(user, k, v)
      
    await session.commit()
    return user
    
    
    
    
    
    
  