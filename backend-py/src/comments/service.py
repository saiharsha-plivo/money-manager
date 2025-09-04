import uuid
from sqlmodel import select
from sqlmodel.ext.asyncio.session import AsyncSession
from fastapi import HTTPException

from src.db.models import Comment, Record
from .schemas import CommentCreate, CommentUpdate


class CommentService:
    async def get_comments_for_record(self, record_id: uuid.UUID, session: AsyncSession):
      statement = select(Comment).where(Comment.record_id == record_id)
      result = await session.exec(statement)
      return result.all()

    async def create_comment_for_record(self,comment_details: CommentCreate,record_id: uuid.UUID,session: AsyncSession) -> Comment:
      record = await session.get(Record, record_id)
      if not record:
        raise HTTPException(status_code=404, detail="Record not found")

      new_comment = Comment(
        **comment_details.model_dump(),
        record_id=record_id,
      )
      session.add(new_comment)
      await session.commit()
      await session.refresh(new_comment)
      return new_comment

    async def update_comment(self,comment_id: uuid.UUID,comment_update: CommentUpdate,session: AsyncSession):
      comment = await session.get(Comment, comment_id)
      if not comment:
          raise HTTPException(status_code=404, detail="Comment not found")

      for key, value in comment_update.model_dump(exclude_unset=True).items():
          setattr(comment, key, value)

      session.add(comment)
      await session.commit()
      await session.refresh(comment)
      return comment

    async def delete_comment(self, comment_id: uuid.UUID, session: AsyncSession):
      comment = await session.get(Comment, comment_id)
      if not comment:
        raise HTTPException(status_code=404, detail="Comment not found")

      await session.delete(comment)
      await session.commit()

      return {"detail": f"Comment {comment_id} deleted successfully"} 