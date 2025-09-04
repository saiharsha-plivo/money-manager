from passlib.context import CryptContext
from src.config import Config
from datetime import datetime, timedelta, timezone
from jose import JWTError, jwt
from src.logger import get_logger
from itsdangerous import URLSafeTimedSerializer

logger = get_logger("auth")

passwd_context = CryptContext(schemes=["bcrypt"])

JWT_SECRET_KEY = Config.JWT_SECRET_KEY
JWT_ALGORITHM = Config.ALGORITHM


def generate_passwd_hash(password: str) -> str:
  hash = passwd_context.hash(password)
  return hash

def verify_password(password: str, hash: str) -> bool:
  return passwd_context.verify(password, hash)

def create_jwt_token(details: dict,expires_delta: timedelta):
  to_encode = details.copy()
  expire = datetime.now(timezone.utc) + expires_delta
  to_encode.update({"exp": expire})
  return jwt.encode(to_encode, JWT_SECRET_KEY, algorithm=JWT_ALGORITHM)
  
def decode_jwt_token(token: str):
  logger.debug(f"token is {token}")
  try:
    payload = jwt.decode(token,JWT_SECRET_KEY,algorithms=[JWT_ALGORITHM])
    return payload.get("sub")
  except JWTError:
    return None

def generate_url_safe_token(email: str):
  serializer = URLSafeTimedSerializer(JWT_SECRET_KEY)
  return serializer.dumps(email)

def decode_url_safe_token(token: str):
  serializer = URLSafeTimedSerializer(JWT_SECRET_KEY)
  try:
    email = serializer.loads(token,max_age=timedelta(days=1).total_seconds())
    return email
  except Exception as e:
    logger.error(f"error verifying url safe token: {e}")
    return None