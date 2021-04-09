package corm

//todo 计算args(参数)数组中的内容长，不计算”“
func Size(args [LEN]string)int{
	i := 0
	for _,s := range args{
		if s == ""{
			continue
		}
		i++
	}
	return i
}

//todo 去除数组中的“”字符
func Remove(args [LEN]string)[LEN]string{
	for i :=range args{
		if i >= len(args){
			break
		}
		if args[i] == ""{
			for j:=i+1;j<len(args);j++{
				args[i]=args[j]
				i++
			}
		}
	}
	return args
}