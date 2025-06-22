class Authorization:
    permissions: dict = {
      "CREATE_MULTIPLE_ACCOUNTS": ["admin", "superuser"],
      "ADD_COMMENT_TO_RECORD": ["admin", "superuser"],
      "EDIT_COMMENT_TO_RECORD": ["admin", "superuser"],
      "DELETE_COMMENT_TO_RECORD": ["admin", "superuser"],
      "GET_COMMENTS_OF_RECORD": ["admin", "superuser"],
    }
    

    def check_access(self, role: str, permission: str) -> bool:
        roles = self.permissions.get(permission, [])
        if role in roles:
            return True
        return False
