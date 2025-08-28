kitex -module iflytek.com/weipan4/learn-go idl/item.thrift
kitex -module iflytek.com/weipan4/learn-go idl/stock.thrift
kitex -module iflytek.com/weipan4/learn-go -service item -use iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen ../../idl/item.thrift
kitex -module iflytek.com/weipan4/learn-go -service stock -use iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen ../../idl/stock.thrift