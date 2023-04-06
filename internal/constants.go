package internal

type TransactionType int

const (
	TRANSACTION_TYPE_OPTIONS TransactionType = iota
	TRANSACTION_TYPE_FOREX
	TRANSACTION_TYPE_DEPOSIT_WITHDRAW
	TRANSACTION_TYPE_BUY_SELL
	TRANSACTION_TYPE_DIVIDEND
	TRANSACTION_TYPE_INTEREST
	TRANSACTION_TYPE_FOREIGN_TAX
)

func (t TransactionType) String() string {
	switch t {
	case TRANSACTION_TYPE_OPTIONS:
		return "options"
	case TRANSACTION_TYPE_FOREX:
		return "forex"
	case TRANSACTION_TYPE_DEPOSIT_WITHDRAW:
		return "deposit-withdraw"
	case TRANSACTION_TYPE_BUY_SELL:
		return "buy-sell"
	case TRANSACTION_TYPE_DIVIDEND:
		return "dividend"
	case TRANSACTION_TYPE_INTEREST:
		return "interest"
	case TRANSACTION_TYPE_FOREIGN_TAX:
		return "foreign-tax"
	default:
		return ""
	}
}

type TransactionsDetailsType int

const (
	TRANSACTIONS_DETAILS_TYPE_DIVIDEND TransactionsDetailsType = iota
	TRANSACTIONS_DETAILS_TYPE_BUY
	TRANSACTIONS_DETAILS_TYPE_SELL
	TRANSACTIONS_DETAILS_TYPE_WITHDRAW
	TRANSACTIONS_DETAILS_TYPE_DEPOSIT
	TRANSACTIONS_DETAILS_TYPE_UNKNOWN
)

func (t TransactionsDetailsType) String() string {
	switch t {
	case TRANSACTIONS_DETAILS_TYPE_DIVIDEND:
		return "DIVIDEND"
	case TRANSACTIONS_DETAILS_TYPE_BUY:
		return "BUY"
	case TRANSACTIONS_DETAILS_TYPE_SELL:
		return "SELL"
	case TRANSACTIONS_DETAILS_TYPE_WITHDRAW:
		return "WITHDRAW"
	case TRANSACTIONS_DETAILS_TYPE_DEPOSIT:
		return "DEPOSIT"
	default:
		return "UNKNOWN"
	}
}

type ChannelType int

const (
	CHANNEL_TYPE_ACCOUNTS ChannelType = iota
	CHANNEL_TYPE_QUOTES
	CHANNEL_TYPE_ORDERDEPTHS
	CHANNEL_TYPE_TRADES
	CHANNEL_TYPE_BROKERTRADESUMMARY
	CHANNEL_TYPE_POSITIONS
	CHANNEL_TYPE_ORDERS
	CHANNEL_TYPE_DEALS
)

func (t ChannelType) String() string {
	switch t {
	case CHANNEL_TYPE_ACCOUNTS:
		return "accounts"
	case CHANNEL_TYPE_QUOTES:
		return "quotes"
	case CHANNEL_TYPE_ORDERDEPTHS:
		return "orderdepths"
	case CHANNEL_TYPE_TRADES:
		return "trades"
	case CHANNEL_TYPE_BROKERTRADESUMMARY:
		return "brokertradesummary"
	case CHANNEL_TYPE_POSITIONS:
		return "positions"
	case CHANNEL_TYPE_ORDERS:
		return "orders"
	case CHANNEL_TYPE_DEALS:
		return "deals"
	default:
		return ""
	}
}

type TimePeriod int

const (
	TIME_PERIOD_TODAY TimePeriod = iota
	TIME_PERIOD_ONE_WEEK
	TIME_PERIOD_ONE_MONTH
	TIME_PERIOD_THREE_MONTHS
	TIME_PERIOD_THIS_YEAR
	TIME_PERIOD_ONE_YEAR
	TIME_PERIOD_THREE_YEARS
	TIME_PERIOD_FIVE_YEARS
	TIME_PERIOD_THREE_YEARS_ROLLING
	TIME_PERIOD_FIVE_YEARS_ROLLING
)

func (t TimePeriod) String() string {
	switch t {
	case TIME_PERIOD_TODAY:
		return "TODAY"
	case TIME_PERIOD_ONE_WEEK:
		return "ONE_WEEK"
	case TIME_PERIOD_ONE_MONTH:
		return "ONE_MONTH"
	case TIME_PERIOD_THREE_MONTHS:
		return "THREE_MONTHS"
	case TIME_PERIOD_THIS_YEAR:
		return "THIS_YEAR"
	case TIME_PERIOD_ONE_YEAR:
		return "ONE_YEAR"
	case TIME_PERIOD_THREE_YEARS:
		return "THREE_YEARS"
	case TIME_PERIOD_FIVE_YEARS:
		return "FIVE_YEARS"
	case TIME_PERIOD_THREE_YEARS_ROLLING:
		return "THREE_YEARS_ROLLING"
	case TIME_PERIOD_FIVE_YEARS_ROLLING:
		return "FIVE_YEARS_ROLLING"
	default:
		return ""
	}
}

type Resolution int

const (
	RESOLUTION_MINUTE Resolution = iota
	RESOLUTION_TWO_MINUTES
	RESOLUTION_FIVE_MINUTES
	RESOLUTION_TEN_MINUTES
	RESOLUTION_THIRTY_MINUTES
	RESOLUTION_HOUR
	RESOLUTION_DAY
	RESOLUTION_WEEK
	RESOLUTION_MONTH
	RESOLUTION_QUARTER
)

func (t Resolution) String() string {
	switch t {
	case RESOLUTION_MINUTE:
		return "MINUTE"
	case RESOLUTION_TWO_MINUTES:
		return "TWO_MINUTES"
	case RESOLUTION_FIVE_MINUTES:
		return "FIVE_MINUTES"
	case RESOLUTION_TEN_MINUTES:
		return "TEN_MINUTES"
	case RESOLUTION_THIRTY_MINUTES:
		return "THIRTY_MINUTES"
	case RESOLUTION_HOUR:
		return "HOUR"
	case RESOLUTION_DAY:
		return "DAY"
	case RESOLUTION_WEEK:
		return "WEEK"
	case RESOLUTION_MONTH:
		return "MONTH"
	case RESOLUTION_QUARTER:
		return "QUARTER"
	default:
		return ""
	}
}

type ListType int

const (
	LIST_TYPE_HIGHEST_RATED_FUNDS ListType = iota
	LIST_TYPE_LOWEST_FEE_INDEX_FUNDS
	LIST_TYPE_BEST_DEVELOPMENT_FUNDS_LAST_THREE_MONTHS
	LIST_TYPE_MOST_OWNED_FUNDS
)

func (t ListType) String() string {
	switch t {
	case LIST_TYPE_HIGHEST_RATED_FUNDS:
		return "HIGHEST_RATED_FUNDS"
	case LIST_TYPE_LOWEST_FEE_INDEX_FUNDS:
		return "LOWEST_FEE_INDEX_FUNDS"
	case LIST_TYPE_BEST_DEVELOPMENT_FUNDS_LAST_THREE_MONTHS:
		return "BEST_DEVELOPMENT_FUNDS_LAST_THREE_MONTHS"
	case LIST_TYPE_MOST_OWNED_FUNDS:
		return "MOST_OWNED_FUNDS"
	default:
		return ""
	}
}

type InstrumentType int

const (
	STOCK InstrumentType = iota
	FUND
	BOND
	OPTION
	FUTURE_FORWARD
	CERTIFICATE
	WARRANT
	EXCHANGE_TRADED_FUND
	INDEX
	PREMIUM_BOND
	SUBSCRIPTION_OPTION
	EQUITY_LINKED_BOND
	CONVERTIBLE
	ANY
)

func (t InstrumentType) String() string {
	switch t {
	case STOCK:
		return "stock"
	case FUND:
		return "fund"
	case BOND:
		return "bond"
	case OPTION:
		return "option"
	case FUTURE_FORWARD:
		return "future_forward"
	case CERTIFICATE:
		return "certificate"
	case WARRANT:
		return "warrant"
	case EXCHANGE_TRADED_FUND:
		return "exchange_traded_fund"
	case INDEX:
		return "index"
	case PREMIUM_BOND:
		return "premium_bond"
	case SUBSCRIPTION_OPTION:
		return "subscription_option"
	case EQUITY_LINKED_BOND:
		return "equity_linked_bond"
	case CONVERTIBLE:
		return "convertible"
	case ANY:
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
	FOLLOW_DOWNWARDS StopLossTriggerType = iota
	FOLLOW_UPWARDS
	LESS_OR_EQUAL
	MORE_OR_EQUAL
)

func (t StopLossTriggerType) String() string {
	switch t {
	case FOLLOW_DOWNWARDS:
		return "FOLLOW_DOWNWARDS"
	case FOLLOW_UPWARDS:
		return "FOLLOW_UPWARDS"
	case LESS_OR_EQUAL:
		return "LESS_OR_EQUAL"
	case MORE_OR_EQUAL:
		return "MORE_OR_EQUAL"
	}
	return ""
}

type StopLossPriceType int

const (
	MONETARY StopLossPriceType = iota
	PERCENTAGE
)

func (t StopLossPriceType) String() string {
	switch t {
	case MONETARY:
		return "MONETARY"
	case PERCENTAGE:
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

package main

import "fmt"

type Route int

const (
	ACCOUNT_OVERVIEW_PATH Route = iota
	ACCOUNTS_POSITIONS_PATH
	AUTHENTICATION_PATH
	CHARTDATA_PATH
	CURRENT_OFFERS_PATH
	DEALS_AND_ORDERS_PATH
	FUND_PATH
	INSIGHTS_PATH
	INSPIRATION_LIST_PATH
	INSTRUMENT_PATH
	INSTRUMENT_DETAILS_PATH
	INSTRUMENT_SEARCH_PATH
	MONTHLY_SAVINGS_CREATE_PATH
	MONTHLY_SAVINGS_PATH
	MONTHLY_SAVINGS_PAUSE_PATH
	MONTHLY_SAVINGS_REMOVE_PATH
	MONTHLY_SAVINGS_RESUME_PATH
	NOTE_PATH
	ORDER_DELETE_PATH
	ORDER_GET_PATH
	ORDER_PLACE_PATH
	ORDER_PLACE_STOP_LOSS_PATH
	ORDER_PLACE_PATH_BUY_FUND
	ORDER_PLACE_PATH_SELL_FUND
	ORDER_EDIT_PATH
	ORDERBOOK_LIST_PATH
	ORDERBOOK_PATH
	OVERVIEW_PATH
	POSITIONS_PATH
	PRICE_ALERT_PATH
	STOP_LOSS_PATH
	TOTP_PATH
	TRANSACTIONS_PATH
	TRANSACTIONS_DETAILS_PATH
	WATCHLISTS_ADD_DELETE_PATH
	WATCHLISTS_PATH
)

func (r Route) String() string {
	switch r {
	case ACCOUNT_OVERVIEW_PATH:
		return "/_mobile/account/{}/overview"
	case ACCOUNTS_POSITIONS_PATH:
		return "/_api/position-data/positions"
	case AUTHENTICATION_PATH:
		return "/_api/authentication/sessions/usercredentials"
	case CHARTDATA_PATH:
		return "/_api/price-chart/stock/{}"
	case CURRENT_OFFERS_PATH:
		return "/_api/customer-offer/currentoffers/"
	case DEALS_AND_ORDERS_PATH:
		return "/_mobile/account/dealsandorders"
	case FUND_PATH:
		return "/_api/fund-guide/guide/{}"
	case INSIGHTS_PATH:
		return "/_api/insights-development/?timePeriod={}&accountIds={}"
	case INSPIRATION_LIST_PATH:
		return "/_mobile/marketing/inspirationlist/{}"
	case INSTRUMENT_PATH:
		return "/_api/market-guide/{}/{}"
	case INSTRUMENT_DETAILS_PATH:
		return "/_api/market-guide/{}/{}/details"
	case INSTRUMENT_SEARCH_PATH:
		return "/_mobile/market/search/{}?query={}&limit={}"
	case MONTHLY_SAVINGS_CREATE_PATH:
		return "/_api/transfer/monthly-savings/{}"
	case MONTHLY_SAVINGS_PATH:
		return "/_mobile/transfer/monthly-savings/{}"
	case MONTHLY_SAVINGS_PAUSE_PATH:
		return "/_api/transfer/monthly-savings/{}/{}/pause"
	case MONTHLY_SAVINGS_REMOVE_PATH:
		return "/_api/transfer/monthly-savings/{}/{}/"
	case MONTHLY_SAVINGS_RESUME_PATH:
		return "/_api/transfer/monthly-savings/{}/{}/resume"
	case NOTE_PATH:
		return "/_api/contract-notes/documents/{}/{}/note.pdf"
	case ORDER_DELETE_PATH:
		return "/_api/order?accountId={}&orderId={}"
	case ORDER_GET_PATH:
		return "/_mobile/order/{}?accountId={}&orderId={}"
	case ORDER_PLACE_PATH:
		return "/_api/trading-critical/rest/order/new"
	case ORDER_PLACE_STOP_LOSS_PATH:
		return "/_api/trading-critical/rest/stoploss/new"
	case ORDER_PLACE_PATH_BUY_FUND:
		return "/_api/fund-guide/fund-order-page/buy"
	case ORDER_PLACE_PATH_SELL_FUND:
		return "/_api/fund-guide/fund-order-page/sell"
	case ORDER_EDIT_PATH:
		return '/_api/order/{}/{}'
	case ORDERBOOK_LIST_PATH:
		return "/_mobile/market/orderbooklist/{}"
	case ORDERBOOK_PATH:
		return "/_mobile/order/{}?orderbookId={}"
	case OVERVIEW_PATH:
		return "/_mobile/account/overview"
	case POSITIONS_PATH:
		return "/_mobile/account/positions"
	case PRICE_ALERT_PATH:
		return "/_cqbe/marketing/service/alert/{}"
	case STOP_LOSS_PATH:
		return "/_api/trading-critical/rest/stoploss"
	case TOTP_PATH:
		return "/_api/authentication/sessions/totp"
	case TRANSACTIONS_PATH:
		return "/_mobile/account/transactions/{}"
	case TRANSACTIONS_DETAILS_PATH:
		return "/_api/transactions"
	case WATCHLISTS_ADD_DELETE_PATH:
		return "/_api/usercontent/watchlist/{}/orderbooks/{}"
	case WATCHLISTS_PATH:
		return "/_mobile/usercontent/watchlist"
	default:
		return ""
	}
}
