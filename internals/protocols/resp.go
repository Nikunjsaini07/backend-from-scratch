package protocols
 
import "errors"

func readLength(data [] byte ) (int , int ){
	 length  := 0  
    
	 for  pos  := range data {
		 if !(data[pos]>='0' && data[pos]<='9'){
			 return length , pos+2
		 }
		 length = length*10 + int(data[pos] - '0')
	 } 

	 return 0 , 0 

}

func readSimpleString(data []byte) (string , int , error){
     pos := 1

	 for ; pos < len(data) && data[pos] != '\r' ; pos++{ 
	 }
	if(pos >= len(data )){
		return "" , 0 , errors.New("invalid RESP")
	 }
	 return string(data[1:pos]) , pos+2 , nil
}

func readError(data []byte) (string , int , error ){
	 return readSimpleString(data)
}
func readInt64(data []byte ) (int64, int , error ){
	pos:=1 
     
	var value int64  = 0

    sign := int64(1)

	if data[pos] == '-' {
		 sign = -1 
		 pos++
	}

	 for ; pos < len(data) && data[pos] != '\r' ; pos++{
		 value = value*10 + int64(data[pos]-'0')
	 }
	 if(pos >= len(data )){
		return 0 , 0 , errors.New("invalid RESP")
	 }

	 return sign*value , pos + 2  , nil  
}

func readBulkstring(data []byte ) ( string , int , error){
	 pos:=1 

	 len , delta := readLength(data[pos:])
	 pos+= delta 


    return string(data[pos:(pos+len)]) , pos+ len + 2 , nil 

}

func readArray(data []byte) (any, int , error) {
	 pos:=1 

	 count , delta := readLength(data[pos:])
	 pos += delta 

	 var elems []any = make([]any , count )
	 for i := range elems {
		 elem , delta , err := DecodeOne(data[pos:])
		 if(err != nil){
			 return nil , 0 , err
		 }
		 elems[i] = elem 
		 pos += delta 
	 }
	 return elems , pos , nil 
}



func DecodeOne(data []byte) (any , int , error){
	if len(data) == 0 {
		return nil , 0 , errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':' :
		return readInt64(data)
    case '$' :
		return readBulkstring(data)
    case '*' : 
	    return readArray(data)
	}

	return nil , 0 ,  errors.New("unknown RESP type")

}

func DecodeArrayString(data []byte ) ([]string , error){
	 value , err := Decode(data )
	 if(err != nil){
		 return nil , err 
	 }
	 ts := value.([]any)
	 tokens := make([]string , len(ts))
	 for i := range tokens {
		 tokens[i] = ts[i].(string)
	 }

	 return tokens , nil 
}


func Decode(data []byte) (any , error){
	if len(data) == 0 {
		 return nil , errors.New("no data")
	}
	value , _, err := DecodeOne(data)
	return value , err 
}