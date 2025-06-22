from sqlmodel.ext.asyncio.session import AsyncSession
from fastapi import HTTPException
from sqlmodel import select
from src.db.models import Account
from .schemas import AccountCreateModel
import uuid
class AccountService:
  async def get_accounts_from_user(self, userid: str, session: AsyncSession):
    statement = (
      select(Account)
      .where(Account.userid == userid)
    )
    result = await session.exec(statement)
    return result.all()
  
  async def create_new_account_for_user(self, account_details: AccountCreateModel, userid: str, session: AsyncSession) -> Account:
    account_details_dict = account_details.model_dump()
    new_account = Account(**account_details_dict)
    new_account.userid = userid
    session.add(new_account)
    await session.commit()
    return new_account
  
  
  async def delete_account(self,account_id: uuid.UUID,session: AsyncSession):
    stmt = select(Account).where(Account.id == account_id)
    account = await session.scalar(stmt)

    if not account:
      raise HTTPException(status_code=404, detail="Account not found")

    await session.delete(account)
    await session.commit()

    return {"detail": f"Account {account_id} deleted successfully"}