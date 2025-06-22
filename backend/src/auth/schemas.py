from pydantic import BaseModel, Field

class UserSignUp(BaseModel):
  username: str = Field(max_length=20,min_length=6)
  password: str = Field(min_length=8)
  email: str 
  
  model_config = {
        "json_schema_extra": {
            "example": {
                "username": "johndoe",
                "email": "johndoe123@co.com",
                "password": "testpass123",
            }
        }
    }


class ForgotPassword(BaseModel):
  password: str = Field(min_length=8)
  confirm_password: str = Field(min_length=8)
  
  model_config = {
    "json_schema_extra": {
      "example": {
        "password": "testpass123",
        "confirm_password": "testpass123",
      }
    }
  }