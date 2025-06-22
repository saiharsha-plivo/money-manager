from fastapi import APIRouter, Depends, HTTPException, status
from src.logger import get_logger
from typing import List
import uuid

from src.db.main import get_session
from sqlmodel.ext.asyncio.session import AsyncSession

from .schemas import Comment, CommentCreate, CommentUpdate
from .service import CommentService
from src.auth.routes import get_current_user
from src.db.models import User
from src.permissions import Authorization

comment_router = APIRouter()
comment_service = CommentService()
authorization = Authorization()
logger = get_logger(__name__)

logger.info("Comment Router Initialized")


@comment_router.get("/records/{record_id}/comments", response_model=List[Comment])
async def get_comments(
  record_id: uuid.UUID,
  session: AsyncSession = Depends(get_session),
  user: User = Depends(get_current_user),
):
  if not authorization.check_access(user.role, "GET_COMMENTS_OF_RECORD"):
    raise HTTPException(
        status_code=status.HTTP_403_FORBIDDEN,
        detail="User does not have permission to get comments",
    )
  return await comment_service.get_comments_for_record(record_id, session)


@comment_router.post("/records/{record_id}/comments", response_model=Comment)
async def create_comment(
  record_id: uuid.UUID,
  comment_details: CommentCreate,
  session: AsyncSession = Depends(get_session),
  user: User = Depends(get_current_user),
):
  if not authorization.check_access(user.role, "ADD_COMMENT_TO_RECORD"):
    raise HTTPException(
        status_code=status.HTTP_403_FORBIDDEN,
        detail="User does not have permission to add comments",
    )
  return await comment_service.create_comment_for_record(
    comment_details, record_id, session
  )


@comment_router.put("/comments/{comment_id}", response_model=Comment)
async def update_comment(
  comment_id: uuid.UUID,
  comment_update: CommentUpdate,
  session: AsyncSession = Depends(get_session),
  user: User = Depends(get_current_user),
):
  if not authorization.check_access(user.role, "EDIT_COMMENT_TO_RECORD"):
    raise HTTPException(
      status_code=status.HTTP_403_FORBIDDEN,
      detail="User does not have permission to edit comments",
    )
  return await comment_service.update_comment(comment_id, comment_update, session)


@comment_router.delete("/comments/{comment_id}", status_code=status.HTTP_200_OK)
async def delete_comment(
    comment_id: uuid.UUID,
    session: AsyncSession = Depends(get_session),
    user: User = Depends(get_current_user),
):
  if not authorization.check_access(user.role, "DELETE_COMMENT_TO_RECORD"):
    raise HTTPException(
      status_code=status.HTTP_403_FORBIDDEN,
      detail="User does not have permission to delete comments",
    )
  return await comment_service.delete_comment(comment_id, session) 