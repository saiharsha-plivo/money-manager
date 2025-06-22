from .config import Config
from typing import List
from fastapi import BackgroundTasks, APIRouter
from fastapi_mail import ConnectionConfig, FastMail, MessageSchema
from pydantic import BaseModel, EmailStr
from starlette.responses import JSONResponse
from .logger import get_logger

logger = get_logger("mail")


conf = ConnectionConfig(
  MAIL_USERNAME=Config.MAIL_USERNAME,
  MAIL_PASSWORD=Config.MAIL_PASSWORD,
  MAIL_FROM=Config.MAIL_FROM,
  MAIL_PORT=Config.MAIL_PORT,
  MAIL_SERVER=Config.MAIL_SERVER,
  MAIL_STARTTLS=Config.MAIL_STARTTLS,
  MAIL_SSL_TLS=Config.MAIL_SSL_TLS,
  USE_CREDENTIALS=True,
  VALIDATE_CERTS=False
)

fm = FastMail(conf)


class EmailSchema(BaseModel):
    email: List[EmailStr]
    message: str
    subject: str


def schedule_email(background_tasks: BackgroundTasks, mail_details: EmailSchema):
  logger.debug(f"Scheduling email to: {mail_details.email} with message: {mail_details.message}")
  message = MessageSchema(
    subject=mail_details.subject,
    recipients=mail_details.email,
    body=mail_details.message,
    subtype="plain",
  )
  background_tasks.add_task(fm.send_message, message)
  logger.info("Email sending task added to background.")


mail_router = APIRouter()


@mail_router.post("/send-email")
async def send_email(background_tasks: BackgroundTasks, mail_details: EmailSchema):
  """
  API Endpoint to send an email.
  """
  schedule_email(background_tasks, mail_details)
  return JSONResponse(status_code=200, content={"message": "email has been sent"})


