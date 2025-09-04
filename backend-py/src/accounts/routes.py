from fastapi import APIRouter, Depends, HTTPException, status
from src.logger import get_logger
from typing import List
from src.db.main import get_session
from sqlmodel.ext.asyncio.session import AsyncSession
from .schemas import Account , AccountCreateModel
from .service import AccountService
from src.auth.routes import get_current_user
from src.db.models import User
from src.permissions import Authorization
import uuid


account_router = APIRouter()
account_service = AccountService()
logger = get_logger(__name__)
authorization = Authorization()

logger.info("Account Router Initialized")


@account_router.get("/", response_model= List[Account])
async def user_accounts(session: AsyncSession = Depends(get_session),user: User = Depends(get_current_user)):
  if user is None:
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED,detail="Issue due invalid token")
  userid = user.id
  logger.debug(f"got a request for account for user {userid}")
  accounts = await account_service.get_accounts_from_user(userid,session)
  logger.debug(f"returned accounts for request for user {userid}")
  return accounts

@account_router.post("/", response_model= Account)
async def create_accounts(account_details: AccountCreateModel, session: AsyncSession = Depends(get_session), user: User = Depends(get_current_user)):
  logger.debug(f"got a request for creating account name {account_details.name} and description is {account_details.description}") 
  if user is None:
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED,detail="Issue due invalid token")
  userid = user.id
  userrole = user.role
  accounts_count = await account_service.get_accounts_from_user(userid,session)
  
  logger.info(f"accounts count: {accounts_count}")
  
  can_create_account = True
  
  if len(accounts_count) > 0 :
    can_create_account = authorization.check_access(userrole,"CREATE_MULTIPLE_ACCOUNTS")  
    
  logger.info(f"can create account : {can_create_account}")
  
  if not can_create_account :
    raise HTTPException(status_code=status.HTTP_405_METHOD_NOT_ALLOWED,detail="Error due to user role (admin , super user can only create multiple accounts)")
  
  account = await account_service.create_new_account_for_user(account_details,userid,session)
  return account


@account_router.delete("/{account_id}",status_code=status.HTTP_200_OK)
async def delete_account(account_id: uuid.UUID,session = Depends(get_session),user: User = Depends(get_current_user)):
  logger.debug(f"got a request to delete account : {account_id}")
  if user is None : 
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED,detail="Issue due invalid token")
  result = await account_service.delete_account(account_id,session)
  return result
  
  
  