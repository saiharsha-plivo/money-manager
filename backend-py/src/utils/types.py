class RecordType:
  types = {
    "1": {
      "type": "income",
      "name": "salary",
    },
    "2": {
      "type": "income",
      "name": "investment",
    },
    "3": {
      "type": "income",
      "name": "part-time job",
    },
    "4": {
      "type": "income",
      "name": "freelance",
    },
    "5": {
      "type": "income",
      "name": "bonus",
    },
    "6": {
      "type": "income",
      "name": "other",
    },
    "7": {
      "type": "expense",
      "name": "food",
    },
    "8": {
      "type": "expense",
      "name": "transportation",
    },
    "9": {
      "type": "expense",
      "name": "shopping",
    },
    "10": {
      "type": "expense",
      "name": "entertainment",
    },
    "11": {
      "type": "expense",
      "name": "health",
    },
    "12": {
      "type": "expense",
      "name": "education",
    },
    "13": {
      "type": "expense",
      "name": "entertainment",    
    },
    "14": {
      "type": "expense",
      "name": "sports",
    },
    "15": {
      "type": "expense",
      "name": "social",
    },
    "16": {
      "type": "expense",
      "name": "addictions",
    },
    "17": {
      "type": "expense",
      "name": "travel",
    },
    "18": {
      "type": "expense",
      "name": "snacks",
    },
    "19": {
      "type": "expense",
      "name": "fruits/vegetables",
    },
    "20": {
      "type": "expense",
      "name": "household",
    },
    "21": {
      "type": "expense",
      "name": "electricity",
    },
    "22": {
      "type": "expense",
      "name": "water",
    },
    "23": {
      "type": "expense",
      "name": "other",
    }
  }
  
  def get_types(self):
    return self.types
  
  def get_type_by_id(self, id):
    return self.types[id]
  
  def get_type_by_name(self, name):
    for type in self.types.values():
      if type["name"] == name:
        return type
    return None
  
  def get_type_by_type(self, type):
    return [type for type in self.types.values() if type["type"] == type]
  