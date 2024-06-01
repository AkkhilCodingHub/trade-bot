import http.client
import json
import backtrader as bt

# Define the trading strategy
class MomentumValueStrategy(bt.Strategy):
    params = (('momentum_period', 30),)  # This defines a tuple of tuples

    def __init__(self):
        # Access the momentum_period correctly using indexing
        momentum_period = self.params.momentum_period  # Corrected access to params using attribute access
        # Initialize the Momentum indicator correctly
        self.momentum = bt.indicators.Momentum(self.datas[0], period=momentum_period)
        self.pe_ratio = self.datas[0].lines.pe  # Assuming PE ratio is part of the data feed

    def next(self):
        if self.momentum[0] > 0 and self.pe_ratio[0] < self.datas[0].lines.pe_avg:
            self.buy(size=100)
            print(f"Bought: {self.datas[0].close[0]} on {self.datas[0].datetime.date(0)}")
        elif self.momentum[0] < 0 or self.pe_ratio[0] > self.datas[0].lines.pe_avg:
            self.sell(size=100)
            print(f"Sold: {self.datas[0].close[0]} on {self.datas[0].datetime.date(0)}")

# Fetch news data
def fetch_news():
    conn = http.client.HTTPSConnection("yahoo-finance15.p.rapidapi.com")
    headers = {
        'X-RapidAPI-Key': "dd43966816msh598c2c63b7c72c1p1782f4jsn940721f52b90",
        'X-RapidAPI-Host': "yahoo-finance15.p.rapidapi.com"
    }
    conn.request("GET", "/api/v2/markets/tickers?type=STOCKS&page=1", headers=headers)
    res = conn.getresponse()
    data = res.read()
    news_data = json.loads(data.decode("utf-8"))
    return news_data

# Backtesting setup
class StockData(bt.feeds.GenericCSVData):
    linesoverride = True  # Corrected by adding linesoverride to allow custom lines
    lines = ('pe', 'pe_avg')
    params = (
        ('dtformat', '%Y-%m-%d'),
        ('datetime', 0),
        ('high', 1),
        ('low', 2),
        ('open', 3),
        ('close', 4),
        ('volume', 5),
        ('openinterest', -1),
        ('pe', 6),
        ('pe_avg', 7)
    )

def main():
    # Create a cerebro entity
    cerebro = bt.Cerebro()

    # Add a strategy
    cerebro.addstrategy(MomentumValueStrategy)

    # Create a data feed
    data = StockData(dataname='stock_data.csv')
    cerebro.adddata(data)

    # Set our desired cash start
    cerebro.broker.set_cash(100000.0)

    # Run over everything
    cerebro.run()

    # Print out the final result
    print(f'Final Portfolio Value: {cerebro.broker.getvalue()}')

if __name__ == "__main__":
    main()

