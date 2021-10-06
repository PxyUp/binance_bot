# Multi-Currency bot for Binance

# Usage

Download from [Release page](https://github.com/PxyUp/binance_bot/releases)

```sh
./$(version)-binance_bot-linux-amd64 --path=config_example.yaml
```

# Important

**If you run not in dryRun mode make sure you have enough money(or other count) for first buy amount of coin which you put in configuration**

# Configuration

[Check example](https://github.com/PxyUp/binance_bot/blob/master/config_example.yaml)

```yaml
binance:
  apiKey: API_KEY # Binance API KEY
  secret: SECRET  # Binance SECRET
  commission: 0.1 # Binance commission in persent
  debug: false # track request to binance api

currencies:
  - symbol: DFUSDT - # Symbol
    # Overwrite global binance config, can be used for another account(section can be removed)
    binance:
      apiKey: API_KEY
      secret: SECRET
      commission: 0.2
      debug: true
    useSR: true # Use Support/Resistance lines (for prevent sell and buy not in best points)
    srDiffPercent: 0.015 # Delta for Support/Resistance price range
    shortPatterns: true # Use short patterns instead of longs
    precision: 4 # How much precision inside price value
    useStopOrder: true # Should we create stop order after we buy
    stopOrderDiff: 0.0030 # If we use stop order, it is difference with real price (for make sure it will not execute immediately)
    count: 1000 # Amount for buy/sell
    countNumber: 0 # For some symbol we cant buy 0.999 for example, and in that case you should put precision number for count here (for example if we can put count 0.888 we should put here 3)
    interval: 60 # How often do action for operation in seconds (for candles 60 is ok, for livePrice mode better use 2sec)
    mode: Candle # We can have Candle/LivePrice
    klineIntervals: # Array of which candles we use for analyse (They analysed together, it is means if 1m show sell and 3m show nothing we will do nothing)
      - 1m
      - 3m
      - 5m
    priceBuffer: 100 # Array of price for analyse (SR/Patterns)
    candleTimeBuffer: 5 # How much we should wait until new candles comming (for example for 1m canldes comming in 00 second, so we will analyse candle on 54second(59-5)
    dryRun: true # Will not do real trade, but we can check profit and estimation of profit)
```