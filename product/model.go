package product

import (
	"bytes"
	"encoding/json"
	"go-api-ws/helpers"
	"io/ioutil"
	"net/http"
)

type SimpleProductStruct struct {
	Index                          string        `json:"_index"`
	Type                           string        `json:"_type"`
	Score                          int           `json:"_score"`
	DocType                        string        `json:"doc_type"`
	ID                             int           `json:"id"`
	Sku                            string        `json:"sku"`
	Name                           string        `json:"name"`
	AttributeSetID                 int           `json:"attribute_set_id"`
	Price                          float64       `json:"price"`
	Status                         int           `json:"status"`
	Visibility                     int           `json:"visibility"`
	TypeID                         string        `json:"type_id"`
	CreatedAt                      string        `json:"created_at"`
	UpdatedAt                      string        `json:"updated_at"`
	CustomAttributes               interface{}   `json:"custom_attributes"`
	FinalPrice                     float64       `json:"final_price"`
	MaxPrice                       float64       `json:"max_price"`
	MaxRegularPrice                float64       `json:"max_regular_price"`
	MinimalRegularPrice            float64       `json:"minimal_regular_price"`
	MinimalPrice                   float64       `json:"minimal_price"`
	RegularPrice                   float64       `json:"regular_price"`
	ItemID                         int           `json:"item_id"`
	ProductID                      int           `json:"product_id"`
	StockID                        int           `json:"stock_id"`
	Qty                            int           `json:"qty"`
	IsInStock                      bool          `json:"is_in_stock"`
	IsQtyDecimal                   bool          `json:"is_qty_decimal"`
	ShowDefaultNotificationMessage bool          `json:"show_default_notification_message"`
	UseConfigMinQty                bool          `json:"use_config_min_qty"`
	MinQty                         int           `json:"min_qty"`
	UseConfigMinSaleQty            int           `json:"use_config_min_sale_qty"`
	MinSaleQty                     int           `json:"min_sale_qty"`
	UseConfigMaxSaleQty            bool          `json:"use_config_max_sale_qty"`
	MaxSaleQty                     int           `json:"max_sale_qty"`
	UseConfigBackorders            bool          `json:"use_config_backorders"`
	Backorders                     int           `json:"backorders"`
	UseConfigNotifyStockQty        bool          `json:"use_config_notify_stock_qty"`
	NotifyStockQty                 int           `json:"notify_stock_qty"`
	UseConfigQtyIncrements         bool          `json:"use_config_qty_increments"`
	QtyIncrements                  int           `json:"qty_increments"`
	UseConfigEnableQtyInc          bool          `json:"use_config_enable_qty_inc"`
	EnableQtyIncrements            bool          `json:"enable_qty_increments"`
	UseConfigManageStock           bool          `json:"use_config_manage_stock"`
	ManageStock                    bool          `json:"manage_stock"`
	LowStockDate                   interface{}   `json:"low_stock_date"`
	IsDecimalDivided               bool          `json:"is_decimal_divided"`
	StockStatusChangedAuto         int           `json:"stock_status_changed_auto"`
	Tsk                            int64         `json:"tsk"`
	Description                    string        `json:"description"`
	Image                          string        `json:"image"`
	SmallImage                     string        `json:"small_image"`
	Thumbnail                      string        `json:"thumbnail"`
	CategoryIds                    []string      `json:"category_ids"`
	OptionsContainer               string        `json:"options_container"`
	RequiredOptions                string        `json:"required_options"`
	HasOptions                     string        `json:"has_options"`
	URLKey                         string        `json:"url_key"`
	TaxClassID                     string        `json:"tax_class_id"`
	Activity                       string        `json:"activity,omitempty"`
	Material                       string        `json:"material,omitempty"`
	Gender                         string        `json:"gender"`
	CategoryGear                   string        `json:"category_gear"`
	ErinRecommends                 string        `json:"erin_recommends"`
	New                            string        `json:"new"`
	ChildDocuments                 []interface{} `json:"_childDocuments_"`
}

type solrResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			JSON string `json:"json"`
		}
	} `json:"responseHeader"`
	Response struct {
		NumFound int                   `json:"numFound"`
		Start    int                   `json:"start"`
		Docs     []SimpleProductStruct `json:"docs"`
	} `json:"response"`
}

func GetProductFromSolrBySKU(sku string) (SimpleProductStruct) {
	request := map[string]interface{}{
		"query": "sku:" + sku,
		"limit": 1}
	requestBytes := new(bytes.Buffer)
	json.NewEncoder(requestBytes).Encode(request)
	resp, err := http.Post(
		"http://localhost:8983/solr/storefrontCore/query",
		"application/json; charset=utf-8",
		requestBytes)
	helpers.PanicErr(err)
	b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%s", b)
	var solrResp solrResponse
	json.Unmarshal(b, &solrResp)
	if solrResp.Response.NumFound > 0 {
		return solrResp.Response.Docs[0]
	}
	return SimpleProductStruct{}
}
