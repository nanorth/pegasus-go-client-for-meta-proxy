// Autogenerated by Thrift Compiler (0.13.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/base"
	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/replication"
	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/rrdb"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var _ = replication.GoUnusedProtection__
var _ = base.GoUnusedProtection__
var _ = rrdb.GoUnusedProtection__

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  update_response put(update_request update)")
	fmt.Fprintln(os.Stderr, "  update_response multi_put(multi_put_request request)")
	fmt.Fprintln(os.Stderr, "  update_response remove(blob key)")
	fmt.Fprintln(os.Stderr, "  multi_remove_response multi_remove(multi_remove_request request)")
	fmt.Fprintln(os.Stderr, "  incr_response incr(incr_request request)")
	fmt.Fprintln(os.Stderr, "  check_and_set_response check_and_set(check_and_set_request request)")
	fmt.Fprintln(os.Stderr, "  check_and_mutate_response check_and_mutate(check_and_mutate_request request)")
	fmt.Fprintln(os.Stderr, "  read_response get(blob key)")
	fmt.Fprintln(os.Stderr, "  multi_get_response multi_get(multi_get_request request)")
	fmt.Fprintln(os.Stderr, "  batch_get_response batch_get(batch_get_request request)")
	fmt.Fprintln(os.Stderr, "  count_response sortkey_count(blob hash_key)")
	fmt.Fprintln(os.Stderr, "  ttl_response ttl(blob key)")
	fmt.Fprintln(os.Stderr, "  scan_response get_scanner(get_scanner_request request)")
	fmt.Fprintln(os.Stderr, "  scan_response scan(scan_request request)")
	fmt.Fprintln(os.Stderr, "  void clear_scanner(i64 context_id)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
	var m map[string]string = h
	return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
	parts := strings.Split(value, ": ")
	if len(parts) != 2 {
		return fmt.Errorf("header should be of format 'Key: Value'")
	}
	h[parts[0]] = parts[1]
	return nil
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	headers := make(httpHeaders)
	var parsedUrl *url.URL
	var trans thrift.TTransport
	_ = strconv.Atoi
	_ = math.Abs
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
	flag.Parse()

	if len(urlString) > 0 {
		var err error
		parsedUrl, err = url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
		if len(headers) > 0 {
			httptrans := trans.(*thrift.THttpClient)
			for key, value := range headers {
				httptrans.SetHeader(key, value)
			}
		}
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	client := rrdb.NewRrdbClient(thrift.NewTStandardClient(iprot, oprot))
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "put":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Put requires 1 args")
			flag.Usage()
		}
		arg39 := flag.Arg(1)
		mbTrans40 := thrift.NewTMemoryBufferLen(len(arg39))
		defer mbTrans40.Close()
		_, err41 := mbTrans40.WriteString(arg39)
		if err41 != nil {
			Usage()
			return
		}
		factory42 := thrift.NewTJSONProtocolFactory()
		jsProt43 := factory42.GetProtocol(mbTrans40)
		argvalue0 := rrdb.NewUpdateRequest()
		err44 := argvalue0.Read(jsProt43)
		if err44 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Put(context.Background(), value0))
		fmt.Print("\n")
		break
	case "multi_put":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "MultiPut requires 1 args")
			flag.Usage()
		}
		arg45 := flag.Arg(1)
		mbTrans46 := thrift.NewTMemoryBufferLen(len(arg45))
		defer mbTrans46.Close()
		_, err47 := mbTrans46.WriteString(arg45)
		if err47 != nil {
			Usage()
			return
		}
		factory48 := thrift.NewTJSONProtocolFactory()
		jsProt49 := factory48.GetProtocol(mbTrans46)
		argvalue0 := rrdb.NewMultiPutRequest()
		err50 := argvalue0.Read(jsProt49)
		if err50 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.MultiPut(context.Background(), value0))
		fmt.Print("\n")
		break
	case "remove":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Remove requires 1 args")
			flag.Usage()
		}
		arg51 := flag.Arg(1)
		mbTrans52 := thrift.NewTMemoryBufferLen(len(arg51))
		defer mbTrans52.Close()
		_, err53 := mbTrans52.WriteString(arg51)
		if err53 != nil {
			Usage()
			return
		}
		factory54 := thrift.NewTJSONProtocolFactory()
		jsProt55 := factory54.GetProtocol(mbTrans52)
		argvalue0 := base.NewBlob()
		err56 := argvalue0.Read(jsProt55)
		if err56 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Remove(context.Background(), value0))
		fmt.Print("\n")
		break
	case "multi_remove":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "MultiRemove requires 1 args")
			flag.Usage()
		}
		arg57 := flag.Arg(1)
		mbTrans58 := thrift.NewTMemoryBufferLen(len(arg57))
		defer mbTrans58.Close()
		_, err59 := mbTrans58.WriteString(arg57)
		if err59 != nil {
			Usage()
			return
		}
		factory60 := thrift.NewTJSONProtocolFactory()
		jsProt61 := factory60.GetProtocol(mbTrans58)
		argvalue0 := rrdb.NewMultiRemoveRequest()
		err62 := argvalue0.Read(jsProt61)
		if err62 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.MultiRemove(context.Background(), value0))
		fmt.Print("\n")
		break
	case "incr":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Incr requires 1 args")
			flag.Usage()
		}
		arg63 := flag.Arg(1)
		mbTrans64 := thrift.NewTMemoryBufferLen(len(arg63))
		defer mbTrans64.Close()
		_, err65 := mbTrans64.WriteString(arg63)
		if err65 != nil {
			Usage()
			return
		}
		factory66 := thrift.NewTJSONProtocolFactory()
		jsProt67 := factory66.GetProtocol(mbTrans64)
		argvalue0 := rrdb.NewIncrRequest()
		err68 := argvalue0.Read(jsProt67)
		if err68 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Incr(context.Background(), value0))
		fmt.Print("\n")
		break
	case "check_and_set":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "CheckAndSet requires 1 args")
			flag.Usage()
		}
		arg69 := flag.Arg(1)
		mbTrans70 := thrift.NewTMemoryBufferLen(len(arg69))
		defer mbTrans70.Close()
		_, err71 := mbTrans70.WriteString(arg69)
		if err71 != nil {
			Usage()
			return
		}
		factory72 := thrift.NewTJSONProtocolFactory()
		jsProt73 := factory72.GetProtocol(mbTrans70)
		argvalue0 := rrdb.NewCheckAndSetRequest()
		err74 := argvalue0.Read(jsProt73)
		if err74 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.CheckAndSet(context.Background(), value0))
		fmt.Print("\n")
		break
	case "check_and_mutate":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "CheckAndMutate requires 1 args")
			flag.Usage()
		}
		arg75 := flag.Arg(1)
		mbTrans76 := thrift.NewTMemoryBufferLen(len(arg75))
		defer mbTrans76.Close()
		_, err77 := mbTrans76.WriteString(arg75)
		if err77 != nil {
			Usage()
			return
		}
		factory78 := thrift.NewTJSONProtocolFactory()
		jsProt79 := factory78.GetProtocol(mbTrans76)
		argvalue0 := rrdb.NewCheckAndMutateRequest()
		err80 := argvalue0.Read(jsProt79)
		if err80 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.CheckAndMutate(context.Background(), value0))
		fmt.Print("\n")
		break
	case "get":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Get requires 1 args")
			flag.Usage()
		}
		arg81 := flag.Arg(1)
		mbTrans82 := thrift.NewTMemoryBufferLen(len(arg81))
		defer mbTrans82.Close()
		_, err83 := mbTrans82.WriteString(arg81)
		if err83 != nil {
			Usage()
			return
		}
		factory84 := thrift.NewTJSONProtocolFactory()
		jsProt85 := factory84.GetProtocol(mbTrans82)
		argvalue0 := base.NewBlob()
		err86 := argvalue0.Read(jsProt85)
		if err86 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Get(context.Background(), value0))
		fmt.Print("\n")
		break
	case "multi_get":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "MultiGet requires 1 args")
			flag.Usage()
		}
		arg87 := flag.Arg(1)
		mbTrans88 := thrift.NewTMemoryBufferLen(len(arg87))
		defer mbTrans88.Close()
		_, err89 := mbTrans88.WriteString(arg87)
		if err89 != nil {
			Usage()
			return
		}
		factory90 := thrift.NewTJSONProtocolFactory()
		jsProt91 := factory90.GetProtocol(mbTrans88)
		argvalue0 := rrdb.NewMultiGetRequest()
		err92 := argvalue0.Read(jsProt91)
		if err92 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.MultiGet(context.Background(), value0))
		fmt.Print("\n")
		break
	case "batch_get":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "BatchGet requires 1 args")
			flag.Usage()
		}
		arg93 := flag.Arg(1)
		mbTrans94 := thrift.NewTMemoryBufferLen(len(arg93))
		defer mbTrans94.Close()
		_, err95 := mbTrans94.WriteString(arg93)
		if err95 != nil {
			Usage()
			return
		}
		factory96 := thrift.NewTJSONProtocolFactory()
		jsProt97 := factory96.GetProtocol(mbTrans94)
		argvalue0 := rrdb.NewBatchGetRequest()
		err98 := argvalue0.Read(jsProt97)
		if err98 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.BatchGet(context.Background(), value0))
		fmt.Print("\n")
		break
	case "sortkey_count":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "SortkeyCount requires 1 args")
			flag.Usage()
		}
		arg99 := flag.Arg(1)
		mbTrans100 := thrift.NewTMemoryBufferLen(len(arg99))
		defer mbTrans100.Close()
		_, err101 := mbTrans100.WriteString(arg99)
		if err101 != nil {
			Usage()
			return
		}
		factory102 := thrift.NewTJSONProtocolFactory()
		jsProt103 := factory102.GetProtocol(mbTrans100)
		argvalue0 := base.NewBlob()
		err104 := argvalue0.Read(jsProt103)
		if err104 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.SortkeyCount(context.Background(), value0))
		fmt.Print("\n")
		break
	case "ttl":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "TTL requires 1 args")
			flag.Usage()
		}
		arg105 := flag.Arg(1)
		mbTrans106 := thrift.NewTMemoryBufferLen(len(arg105))
		defer mbTrans106.Close()
		_, err107 := mbTrans106.WriteString(arg105)
		if err107 != nil {
			Usage()
			return
		}
		factory108 := thrift.NewTJSONProtocolFactory()
		jsProt109 := factory108.GetProtocol(mbTrans106)
		argvalue0 := base.NewBlob()
		err110 := argvalue0.Read(jsProt109)
		if err110 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.TTL(context.Background(), value0))
		fmt.Print("\n")
		break
	case "get_scanner":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "GetScanner requires 1 args")
			flag.Usage()
		}
		arg111 := flag.Arg(1)
		mbTrans112 := thrift.NewTMemoryBufferLen(len(arg111))
		defer mbTrans112.Close()
		_, err113 := mbTrans112.WriteString(arg111)
		if err113 != nil {
			Usage()
			return
		}
		factory114 := thrift.NewTJSONProtocolFactory()
		jsProt115 := factory114.GetProtocol(mbTrans112)
		argvalue0 := rrdb.NewGetScannerRequest()
		err116 := argvalue0.Read(jsProt115)
		if err116 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.GetScanner(context.Background(), value0))
		fmt.Print("\n")
		break
	case "scan":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Scan requires 1 args")
			flag.Usage()
		}
		arg117 := flag.Arg(1)
		mbTrans118 := thrift.NewTMemoryBufferLen(len(arg117))
		defer mbTrans118.Close()
		_, err119 := mbTrans118.WriteString(arg117)
		if err119 != nil {
			Usage()
			return
		}
		factory120 := thrift.NewTJSONProtocolFactory()
		jsProt121 := factory120.GetProtocol(mbTrans118)
		argvalue0 := rrdb.NewScanRequest()
		err122 := argvalue0.Read(jsProt121)
		if err122 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Scan(context.Background(), value0))
		fmt.Print("\n")
		break
	case "clear_scanner":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ClearScanner requires 1 args")
			flag.Usage()
		}
		argvalue0, err123 := (strconv.ParseInt(flag.Arg(1), 10, 64))
		if err123 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ClearScanner(context.Background(), value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
