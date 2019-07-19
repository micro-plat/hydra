package sqls

const Get_ORDER_ID = `select seq_order_no.nextval from dual`
const ORDER_QUERY = "SELECT t.order_id,t.merchant_id,t.merchant_order_no,t.account,t.status FROM ORDER_MAIN T WHERE T.MERCHANT_ID=@merchant_id and t.merchant_order_no=@order_no"
const ORDER_CREATE = "INSERT INTO ORDER_MAIN(MERCHANT_ID,MERCHANT_ORDER_NO,ACCOUNT,FACE,NUM)vallues(@merchant_id,@order_no,@account,@face,@num)"
