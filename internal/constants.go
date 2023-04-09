package internal

type TransactionType int

const (
	Options TransactionType = iota
	Forex
	DepositWithdraw
	BuySell
	DividendTransaction
	Interest
	ForeignTax
)

func (t TransactionType) String() string {
	switch t {
	case Options:
		return "options"
	case Forex:
		return "forex"
	case DepositWithdraw:
		return "deposit-withdraw"
	case BuySell:
		return "buy-sell"
	case DividendTransaction:
		return "dividend"
	case Interest:
		return "interest"
	case ForeignTax:
		return "foreign-tax"
	default:
		return ""
	}
}

type TransactionsDetailsType int

const (
	Dividend TransactionsDetailsType = iota
	Buy
	Sell
	Withdraw
	Deposit
)

func (t TransactionsDetailsType) String() string {
	switch t {
	case Dividend:
		return "DIVIDEND"
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	case Withdraw:
		return "WITHDRAW"
	case Deposit:
		return "DEPOSIT"
	default:
		return "UNKNOWN"
	}
}

type ChannelType int

const (
	Accounts ChannelType = iota
	Quotes
	OrderDepths
	Trades
	BrokerTradeSummary
	Positions
	Orders
	TypeDeals
)

func (t ChannelType) String() string {
	switch t {
	case Accounts:
		return "accounts"
	case Quotes:
		return "quotes"
	case OrderDepths:
		return "orderdepths"
	case Trades:
		return "trades"
	case BrokerTradeSummary:
		return "brokertradesummary"
	case Positions:
		return "positions"
	case Orders:
		return "orders"
	case TypeDeals:
		return "deals"
	default:
		return ""
	}
}

type TimePeriod int

const (
	Today TimePeriod = iota
	OneWeek
	OneMonth
	ThreeMonths
	ThisYear
	OneYear
	ThreeYears
	FiveYears
	ThreeYearsRolling
	FiveYearsRolling
)

func (t TimePeriod) String() string {
	switch t {
	case Today:
		return "TODAY"
	case OneWeek:
		return "ONE_WEEK"
	case OneMonth:
		return "ONE_MONTH"
	case ThreeMonths:
		return "THREE_MONTHS"
	case ThisYear:
		return "THIS_YEAR"
	case OneYear:
		return "ONE_YEAR"
	case ThreeYears:
		return "THREE_YEARS"
	case FiveYears:
		return "FIVE_YEARS"
	case ThreeYearsRolling:
		return "THREE_YEARS_ROLLING"
	case FiveYearsRolling:
		return "FIVE_YEARS_ROLLING"
	default:
		return ""
	}
}

type Resolution int

const (
	Minute Resolution = iota
	TwoMinutes
	FiveMinutes
	TenMinutes
	ThirtyMinutes
	Hour
	Day
	Week
	Month
	Quarter
)

func (t Resolution) String() string {
	switch t {
	case Minute:
		return "MINUTE"
	case TwoMinutes:
		return "TWO_MINUTES"
	case FiveMinutes:
		return "FIVE_MINUTES"
	case TenMinutes:
		return "TEN_MINUTES"
	case ThirtyMinutes:
		return "THIRTY_MINUTES"
	case Hour:
		return "HOUR"
	case Day:
		return "DAY"
	case Week:
		return "WEEK"
	case Month:
		return "MONTH"
	case Quarter:
		return "QUARTER"
	default:
		return ""
	}
}

type ListType int

const (
	HighestRatedFunds ListType = iota
	LowestFeeIndexFunds
	BestDevelopmentFundsLastThreeMonths
	MostOwnedFunds
)

func (t ListType) String() string {
	switch t {
	case HighestRatedFunds:
		return "HIGHEST_RATED_FUNDS"
	case LowestFeeIndexFunds:
		return "LOWEST_FEE_INDEX_FUNDS"
	case BestDevelopmentFundsLastThreeMonths:
		return "BEST_DEVELOPMENT_FUNDS_LAST_THREE_MONTHS"
	case MostOwnedFunds:
		return "MOST_OWNED_FUNDS"
	default:
		return ""
	}
}

type InstrumentType int

const (
	Stock InstrumentType = iota
	Fund
	Bond
	Option
	FutureForward
	Certificate
	Warrant
	ExchangeTradedFund
	Index
	PremiumBond
	SubscriptionOption
	EquityLinkedBond
	Convertible
	Any
)

func (t InstrumentType) String() string {
	switch t {
	case Stock:
		return "stock"
	case Fund:
		return "fund"
	case Bond:
		return "bond"
	case Option:
		return "option"
	case FutureForward:
		return "future_forward"
	case Certificate:
		return "certificate"
	case Warrant:
		return "warrant"
	case ExchangeTradedFund:
		return "exchange_traded_fund"
	case Index:
		return "index"
	case PremiumBond:
		return "premium_bond"
	case SubscriptionOption:
		return "subscription_option"
	case EquityLinkedBond:
		return "equity_linked_bond"
	case Convertible:
		return "convertible"
	case Any:
		return ""
	}
	return ""
}

type OrderType int

const (
	BUY OrderType = iota
	SELL
)

func (t OrderType) String() string {
	switch t {
	case BUY:
		return "BUY"
	case SELL:
		return "SELL"
	}
	return ""
}

type StopLossTriggerType int

const (
	FollowDownwards StopLossTriggerType = iota
	FollowUpwards
	LessOrEqual
	MoreOrEqual
)

func (t StopLossTriggerType) String() string {
	switch t {
	case FollowDownwards:
		return "FOLLOW_DOWNWARDS"
	case FollowUpwards:
		return "FOLLOW_UPWARDS"
	case LessOrEqual:
		return "LESS_OR_EQUAL"
	case MoreOrEqual:
		return "MORE_OR_EQUAL"
	}
	return ""
}

type StopLossPriceType int

const (
	Monetary StopLossPriceType = iota
	Percentage
)

func (t StopLossPriceType) String() string {
	switch t {
	case Monetary:
		return "MONETARY"
	case Percentage:
		return "PERCENTAGE"
	}
	return ""
}

type HttpMethod int

const (
	POST HttpMethod = iota + 1
	GET
	PUT
	DELETE
)

func (t HttpMethod) String() string {
	switch t {
	case POST:
		return "POST"
	case GET:
		return "GET"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	}
	return ""
}

type Route int

const (
	AccountOverviewPath Route = iota
	AccountsPositionsPath
	AuthenticationPath
	ChartdataPath
	CurrentOffersPath
	DealsAndOrdersPath
	FundPath
	InsightsPath
	InspirationListPath
	InstrumentPath
	InstrumentDetailsPath
	InstrumentSearchPath
	MonthlySavingsCreatePath
	MonthlySavingsPath
	MonthlySavingsPausePath
	MonthlySavingsRemovePath
	MonthlySavingsResumePath
	NotePath
	OrderDeletePath
	OrderGetPath
	OrderPlacePath
	OrderPlaceStopLossPath
	OrderPlacePathBuyFund
	OrderPlacePathSellFund
	OrderEditPath
	OrderbookListPath
	OrderbookPath
	OverviewPath
	PositionsPath
	PriceAlertPath
	StopLossPath
	TotpPath
	TransactionsPath
	TransactionsDetailsPath
	WatchlistsAddDeletePath
	WatchlistsPath
)

func (r Route) String() string {
	switch r {
	case AccountOverviewPath:
		return "/_mobile/account/{}/overview"
	case AccountsPositionsPath:
		return "/_api/position-data/positions"
	case AuthenticationPath:
		return "/_api/authentication/sessions/usercredentials"
	case ChartdataPath:
		return "/_api/price-chart/stock/{}"
	case CurrentOffersPath:
		return "/_api/customer-offer/currentoffers/"
	case DealsAndOrdersPath:
		return "/_mobile/account/dealsandorders"
	case FundPath:
		return "/_api/fund-guide/guide/{}"
	case InsightsPath:
		return "/_api/insights-development/?timePeriod={}&accountIds={}"
	case InspirationListPath:
		return "/_mobile/marketing/inspirationlist/{}"
	case InstrumentPath:
		return "/_api/market-guide/{}/{}"
	case InstrumentDetailsPath:
		return "/_api/market-guide/{}/{}/details"
	case InstrumentSearchPath:
		return "/_mobile/market/search/{}?query={}&limit={}"
	case MonthlySavingsCreatePath:
		return "/_api/transfer/monthly-savings/{}"
	case MonthlySavingsPath:
		return "/_mobile/transfer/monthly-savings/{}"
	case MonthlySavingsPausePath:
		return "/_api/transfer/monthly-savings/{}/{}/pause"
	case MonthlySavingsRemovePath:
		return "/_api/transfer/monthly-savings/{}/{}/"
	case MonthlySavingsResumePath:
		return "/_api/transfer/monthly-savings/{}/{}/resume"
	case NotePath:
		return "/_api/contract-notes/documents/{}/{}/note.pdf"
	case OrderDeletePath:
		return "/_api/order?accountId={}&orderId={}"
	case OrderGetPath:
		return "/_mobile/order/{}?accountId={}&orderId={}"
	case OrderPlacePath:
		return "/_api/trading-critical/rest/order/new"
	case OrderPlaceStopLossPath:
		return "/_api/trading-critical/rest/stoploss/new"
	case OrderPlacePathBuyFund:
		return "/_api/fund-guide/fund-order-page/buy"
	case OrderPlacePathSellFund:
		return "/_api/fund-guide/fund-order-page/sell"
	case OrderEditPath:
		return "/_api/order/{}/{}"
	case OrderbookListPath:
		return "/_mobile/market/orderbooklist/{}"
	case OrderbookPath:
		return "/_mobile/order/{}?orderbookId={}"
	case OverviewPath:
		return "/_mobile/account/overview"
	case PositionsPath:
		return "/_mobile/account/positions"
	case PriceAlertPath:
		return "/_cqbe/marketing/service/alert/{}"
	case StopLossPath:
		return "/_api/trading-critical/rest/stoploss"
	case TotpPath:
		return "/_api/authentication/sessions/totp"
	case TransactionsPath:
		return "/_mobile/account/transactions/{}"
	case TransactionsDetailsPath:
		return "/_api/transactions"
	case WatchlistsAddDeletePath:
		return "/_api/usercontent/watchlist/{}/orderbooks/{}"
	case WatchlistsPath:
		return "/_mobile/usercontent/watchlist"
	default:
		return ""
	}
}
