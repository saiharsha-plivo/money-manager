from fastapi import FastAPI
from .logger import get_logger
from contextlib import asynccontextmanager
from .db.main import init_db
from .accounts.routes import account_router
from .auth.routes import auth_router
from .records.routes import record_router
from .mail import mail_router
from .comments.routes import comment_router


version = "v1"

description = """
  A rest api for money manager , it helps to authenticate 
  authorize and manage the finances of users 
  """

version_prefix = f"/api/{version}"

logger = get_logger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    await init_db()
    logger.info("database has initialized")
    yield


app = FastAPI(
    title="Money-Manager",
    description=description,
    version=version,
    lifespan=lifespan,
    contact={
        "name": "Sai Harsha",
        "url": "https://github.com/harsha190202",
        "email": "harshanaidu3330@gmail.com",
    },
    openapi_url=f"{version_prefix}/openapi.json",
    docs_url=f"{version_prefix}/docs",
    redoc_url=f"{version_prefix}/redoc",
)

logger.info("Application has started, FastAPI initialized")


@app.get("/")
async def home():
    return {"message": "Money Manager Home route"}

@app.get("/health-check")
async def health_check():
  return {"status": "healthy"}

app.include_router(account_router, prefix=f"{version_prefix}/accounts", tags=["accounts"])
app.include_router(auth_router, prefix=f"{version_prefix}/auth", tags=["auth,users"])
app.include_router(record_router, prefix=f"{version_prefix}/accounts", tags=["records"])
app.include_router(mail_router, prefix=f"{version_prefix}/mail", tags=["mail"])
app.include_router(comment_router, prefix=f"{version_prefix}", tags=["comments"])