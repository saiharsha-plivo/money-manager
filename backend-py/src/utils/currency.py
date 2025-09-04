class Currency:
  currency_list = {
    "1": {"name": "US Dollar", "symbol": "$", "code": "USD", "rate": 1.00},
    "2": {"name": "Euro", "symbol": "€", "code": "EUR", "rate": 0.85},
    "3": {"name": "British Pound", "symbol": "£", "code": "GBP", "rate": 0.75},
    "4": {"name": "Japanese Yen", "symbol": "¥", "code": "JPY", "rate": 110.0},
    "5": {"name": "Swiss Franc", "symbol": "CHF", "code": "CHF", "rate": 0.91},
    "6": {"name": "Canadian Dollar", "symbol": "C$", "code": "CAD", "rate": 1.26},
    "7": {"name": "Australian Dollar", "symbol": "A$", "code": "AUD", "rate": 1.34},
    "8": {"name": "Indian Rupee", "symbol": "₹", "code": "INR", "rate": 74.5},
    "9": {"name": "Chinese Yuan", "symbol": "¥", "code": "CNY", "rate": 6.45},
    "10": {"name": "Singapore Dollar", "symbol": "S$", "code": "SGD", "rate": 1.35},
    "11": {"name": "New Zealand Dollar", "symbol": "NZ$", "code": "NZD", "rate": 1.43},
    "12": {"name": "South Korean Won", "symbol": "₩", "code": "KRW", "rate": 1150.0},
    "13": {"name": "Mexican Peso", "symbol": "Mex$", "code": "MXN", "rate": 20.0},
    "14": {"name": "Brazilian Real", "symbol": "R$", "code": "BRL", "rate": 5.2},
    "15": {"name": "South African Rand", "symbol": "R", "code": "ZAR", "rate": 14.3},
    "16": {"name": "Russian Ruble", "symbol": "₽", "code": "RUB", "rate": 73.0},
    "17": {"name": "Turkish Lira", "symbol": "₺", "code": "TRY", "rate": 8.6},
    "18": {"name": "Swedish Krona", "symbol": "kr", "code": "SEK", "rate": 8.5},
    "19": {"name": "Norwegian Krone", "symbol": "kr", "code": "NOK", "rate": 8.7},
    "20": {"name": "Danish Krone", "symbol": "kr", "code": "DKK", "rate": 6.3},
    "21": {"name": "Polish Zloty", "symbol": "zł", "code": "PLN", "rate": 3.9},
    "22": {"name": "Thai Baht", "symbol": "฿", "code": "THB", "rate": 33.0},
    "23": {"name": "Malaysian Ringgit", "symbol": "RM", "code": "MYR", "rate": 4.1},
    "24": {"name": "Indonesian Rupiah", "symbol": "Rp", "code": "IDR", "rate": 14200.0},
    "25": {"name": "Philippine Peso", "symbol": "₱", "code": "PHP", "rate": 50.0},
    "26": {"name": "Vietnamese Dong", "symbol": "₫", "code": "VND", "rate": 23000.0},
    "27": {"name": "UAE Dirham", "symbol": "د.إ", "code": "AED", "rate": 3.67},
    "28": {"name": "Saudi Riyal", "symbol": "﷼", "code": "SAR", "rate": 3.75},
    "29": {"name": "Israeli Shekel", "symbol": "₪", "code": "ILS", "rate": 3.3},
    "30": {"name": "Argentine Peso", "symbol": "$", "code": "ARS", "rate": 95.0}
  }
  
  def get_currency_list(self):
    return self.currency_list
  
  def get_currency_by_id(self, id):
    return self.currency_list[id]
  
  def usd_to_currency(self, amount, currency_id):
    currency = self.get_currency_by_id(currency_id)
    return amount * currency["rate"]
  
  def currency_to_usd(self, amount, currency_id):
    currency = self.get_currency_by_id(currency_id)
    return amount / currency["rate"]
  
