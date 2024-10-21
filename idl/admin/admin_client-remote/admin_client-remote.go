// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/incubator-pegasus/idl/admin"
	"github.com/apache/incubator-pegasus/idl/base"
	"github.com/apache/incubator-pegasus/idl/replication"
	"github.com/apache/thrift/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var _ = base.GoUnusedProtection__
var _ = replication.GoUnusedProtection__

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  configuration_create_app_response create_app(configuration_create_app_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_drop_app_response drop_app(configuration_drop_app_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_recall_app_response recall_app(configuration_recall_app_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_list_apps_response list_apps(configuration_list_apps_request req)")
	fmt.Fprintln(os.Stderr, "  duplication_add_response add_duplication(duplication_add_request req)")
	fmt.Fprintln(os.Stderr, "  duplication_query_response query_duplication(duplication_query_request req)")
	fmt.Fprintln(os.Stderr, "  duplication_modify_response modify_duplication(duplication_modify_request req)")
	fmt.Fprintln(os.Stderr, "  query_app_info_response query_app_info(query_app_info_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_update_app_env_response update_app_env(configuration_update_app_env_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_list_nodes_response list_nodes(configuration_list_nodes_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_cluster_info_response query_cluster_info(configuration_cluster_info_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_meta_control_response meta_control(configuration_meta_control_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_query_backup_policy_response query_backup_policy(configuration_query_backup_policy_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_balancer_response balance(configuration_balancer_request req)")
	fmt.Fprintln(os.Stderr, "  start_backup_app_response start_backup_app(start_backup_app_request req)")
	fmt.Fprintln(os.Stderr, "  query_backup_status_response query_backup_status(query_backup_status_request req)")
	fmt.Fprintln(os.Stderr, "  configuration_create_app_response restore_app(configuration_restore_request req)")
	fmt.Fprintln(os.Stderr, "  start_partition_split_response start_partition_split(start_partition_split_request req)")
	fmt.Fprintln(os.Stderr, "  query_split_response query_split_status(query_split_request req)")
	fmt.Fprintln(os.Stderr, "  control_split_response control_partition_split(control_split_request req)")
	fmt.Fprintln(os.Stderr, "  start_bulk_load_response start_bulk_load(start_bulk_load_request req)")
	fmt.Fprintln(os.Stderr, "  query_bulk_load_response query_bulk_load_status(query_bulk_load_request req)")
	fmt.Fprintln(os.Stderr, "  control_bulk_load_response control_bulk_load(control_bulk_load_request req)")
	fmt.Fprintln(os.Stderr, "  clear_bulk_load_state_response clear_bulk_load(clear_bulk_load_state_request req)")
	fmt.Fprintln(os.Stderr, "  start_app_manual_compact_response start_manual_compact(start_app_manual_compact_request req)")
	fmt.Fprintln(os.Stderr, "  query_app_manual_compact_response query_manual_compact(query_app_manual_compact_request req)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
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
	flag.Parse()

	if len(urlString) > 0 {
		var err error
		parsedUrl, err = url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
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
	client := admin.NewAdminClientClient(thrift.NewTStandardClient(iprot, oprot))
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "create_app":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "CreateApp requires 1 args")
			flag.Usage()
		}
		arg73 := flag.Arg(1)
		mbTrans74 := thrift.NewTMemoryBufferLen(len(arg73))
		defer mbTrans74.Close()
		_, err75 := mbTrans74.WriteString(arg73)
		if err75 != nil {
			Usage()
			return
		}
		factory76 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt77 := factory76.GetProtocol(mbTrans74)
		argvalue0 := admin.NewConfigurationCreateAppRequest()
		err78 := argvalue0.Read(jsProt77)
		if err78 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.CreateApp(context.Background(), value0))
		fmt.Print("\n")
		break
	case "drop_app":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "DropApp requires 1 args")
			flag.Usage()
		}
		arg79 := flag.Arg(1)
		mbTrans80 := thrift.NewTMemoryBufferLen(len(arg79))
		defer mbTrans80.Close()
		_, err81 := mbTrans80.WriteString(arg79)
		if err81 != nil {
			Usage()
			return
		}
		factory82 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt83 := factory82.GetProtocol(mbTrans80)
		argvalue0 := admin.NewConfigurationDropAppRequest()
		err84 := argvalue0.Read(jsProt83)
		if err84 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.DropApp(context.Background(), value0))
		fmt.Print("\n")
		break
	case "recall_app":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "RecallApp requires 1 args")
			flag.Usage()
		}
		arg85 := flag.Arg(1)
		mbTrans86 := thrift.NewTMemoryBufferLen(len(arg85))
		defer mbTrans86.Close()
		_, err87 := mbTrans86.WriteString(arg85)
		if err87 != nil {
			Usage()
			return
		}
		factory88 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt89 := factory88.GetProtocol(mbTrans86)
		argvalue0 := admin.NewConfigurationRecallAppRequest()
		err90 := argvalue0.Read(jsProt89)
		if err90 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.RecallApp(context.Background(), value0))
		fmt.Print("\n")
		break
	case "list_apps":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ListApps requires 1 args")
			flag.Usage()
		}
		arg91 := flag.Arg(1)
		mbTrans92 := thrift.NewTMemoryBufferLen(len(arg91))
		defer mbTrans92.Close()
		_, err93 := mbTrans92.WriteString(arg91)
		if err93 != nil {
			Usage()
			return
		}
		factory94 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt95 := factory94.GetProtocol(mbTrans92)
		argvalue0 := admin.NewConfigurationListAppsRequest()
		err96 := argvalue0.Read(jsProt95)
		if err96 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ListApps(context.Background(), value0))
		fmt.Print("\n")
		break
	case "add_duplication":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "AddDuplication requires 1 args")
			flag.Usage()
		}
		arg97 := flag.Arg(1)
		mbTrans98 := thrift.NewTMemoryBufferLen(len(arg97))
		defer mbTrans98.Close()
		_, err99 := mbTrans98.WriteString(arg97)
		if err99 != nil {
			Usage()
			return
		}
		factory100 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt101 := factory100.GetProtocol(mbTrans98)
		argvalue0 := admin.NewDuplicationAddRequest()
		err102 := argvalue0.Read(jsProt101)
		if err102 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.AddDuplication(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_duplication":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryDuplication requires 1 args")
			flag.Usage()
		}
		arg103 := flag.Arg(1)
		mbTrans104 := thrift.NewTMemoryBufferLen(len(arg103))
		defer mbTrans104.Close()
		_, err105 := mbTrans104.WriteString(arg103)
		if err105 != nil {
			Usage()
			return
		}
		factory106 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt107 := factory106.GetProtocol(mbTrans104)
		argvalue0 := admin.NewDuplicationQueryRequest()
		err108 := argvalue0.Read(jsProt107)
		if err108 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryDuplication(context.Background(), value0))
		fmt.Print("\n")
		break
	case "modify_duplication":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ModifyDuplication requires 1 args")
			flag.Usage()
		}
		arg109 := flag.Arg(1)
		mbTrans110 := thrift.NewTMemoryBufferLen(len(arg109))
		defer mbTrans110.Close()
		_, err111 := mbTrans110.WriteString(arg109)
		if err111 != nil {
			Usage()
			return
		}
		factory112 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt113 := factory112.GetProtocol(mbTrans110)
		argvalue0 := admin.NewDuplicationModifyRequest()
		err114 := argvalue0.Read(jsProt113)
		if err114 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ModifyDuplication(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_app_info":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryAppInfo requires 1 args")
			flag.Usage()
		}
		arg115 := flag.Arg(1)
		mbTrans116 := thrift.NewTMemoryBufferLen(len(arg115))
		defer mbTrans116.Close()
		_, err117 := mbTrans116.WriteString(arg115)
		if err117 != nil {
			Usage()
			return
		}
		factory118 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt119 := factory118.GetProtocol(mbTrans116)
		argvalue0 := admin.NewQueryAppInfoRequest()
		err120 := argvalue0.Read(jsProt119)
		if err120 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryAppInfo(context.Background(), value0))
		fmt.Print("\n")
		break
	case "update_app_env":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "UpdateAppEnv requires 1 args")
			flag.Usage()
		}
		arg121 := flag.Arg(1)
		mbTrans122 := thrift.NewTMemoryBufferLen(len(arg121))
		defer mbTrans122.Close()
		_, err123 := mbTrans122.WriteString(arg121)
		if err123 != nil {
			Usage()
			return
		}
		factory124 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt125 := factory124.GetProtocol(mbTrans122)
		argvalue0 := admin.NewConfigurationUpdateAppEnvRequest()
		err126 := argvalue0.Read(jsProt125)
		if err126 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.UpdateAppEnv(context.Background(), value0))
		fmt.Print("\n")
		break
	case "list_nodes":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ListNodes requires 1 args")
			flag.Usage()
		}
		arg127 := flag.Arg(1)
		mbTrans128 := thrift.NewTMemoryBufferLen(len(arg127))
		defer mbTrans128.Close()
		_, err129 := mbTrans128.WriteString(arg127)
		if err129 != nil {
			Usage()
			return
		}
		factory130 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt131 := factory130.GetProtocol(mbTrans128)
		argvalue0 := admin.NewConfigurationListNodesRequest()
		err132 := argvalue0.Read(jsProt131)
		if err132 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ListNodes(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_cluster_info":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryClusterInfo requires 1 args")
			flag.Usage()
		}
		arg133 := flag.Arg(1)
		mbTrans134 := thrift.NewTMemoryBufferLen(len(arg133))
		defer mbTrans134.Close()
		_, err135 := mbTrans134.WriteString(arg133)
		if err135 != nil {
			Usage()
			return
		}
		factory136 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt137 := factory136.GetProtocol(mbTrans134)
		argvalue0 := admin.NewConfigurationClusterInfoRequest()
		err138 := argvalue0.Read(jsProt137)
		if err138 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryClusterInfo(context.Background(), value0))
		fmt.Print("\n")
		break
	case "meta_control":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "MetaControl requires 1 args")
			flag.Usage()
		}
		arg139 := flag.Arg(1)
		mbTrans140 := thrift.NewTMemoryBufferLen(len(arg139))
		defer mbTrans140.Close()
		_, err141 := mbTrans140.WriteString(arg139)
		if err141 != nil {
			Usage()
			return
		}
		factory142 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt143 := factory142.GetProtocol(mbTrans140)
		argvalue0 := admin.NewConfigurationMetaControlRequest()
		err144 := argvalue0.Read(jsProt143)
		if err144 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.MetaControl(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_backup_policy":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryBackupPolicy requires 1 args")
			flag.Usage()
		}
		arg145 := flag.Arg(1)
		mbTrans146 := thrift.NewTMemoryBufferLen(len(arg145))
		defer mbTrans146.Close()
		_, err147 := mbTrans146.WriteString(arg145)
		if err147 != nil {
			Usage()
			return
		}
		factory148 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt149 := factory148.GetProtocol(mbTrans146)
		argvalue0 := admin.NewConfigurationQueryBackupPolicyRequest()
		err150 := argvalue0.Read(jsProt149)
		if err150 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryBackupPolicy(context.Background(), value0))
		fmt.Print("\n")
		break
	case "balance":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Balance requires 1 args")
			flag.Usage()
		}
		arg151 := flag.Arg(1)
		mbTrans152 := thrift.NewTMemoryBufferLen(len(arg151))
		defer mbTrans152.Close()
		_, err153 := mbTrans152.WriteString(arg151)
		if err153 != nil {
			Usage()
			return
		}
		factory154 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt155 := factory154.GetProtocol(mbTrans152)
		argvalue0 := admin.NewConfigurationBalancerRequest()
		err156 := argvalue0.Read(jsProt155)
		if err156 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Balance(context.Background(), value0))
		fmt.Print("\n")
		break
	case "start_backup_app":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "StartBackupApp requires 1 args")
			flag.Usage()
		}
		arg157 := flag.Arg(1)
		mbTrans158 := thrift.NewTMemoryBufferLen(len(arg157))
		defer mbTrans158.Close()
		_, err159 := mbTrans158.WriteString(arg157)
		if err159 != nil {
			Usage()
			return
		}
		factory160 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt161 := factory160.GetProtocol(mbTrans158)
		argvalue0 := admin.NewStartBackupAppRequest()
		err162 := argvalue0.Read(jsProt161)
		if err162 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.StartBackupApp(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_backup_status":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryBackupStatus requires 1 args")
			flag.Usage()
		}
		arg163 := flag.Arg(1)
		mbTrans164 := thrift.NewTMemoryBufferLen(len(arg163))
		defer mbTrans164.Close()
		_, err165 := mbTrans164.WriteString(arg163)
		if err165 != nil {
			Usage()
			return
		}
		factory166 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt167 := factory166.GetProtocol(mbTrans164)
		argvalue0 := admin.NewQueryBackupStatusRequest()
		err168 := argvalue0.Read(jsProt167)
		if err168 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryBackupStatus(context.Background(), value0))
		fmt.Print("\n")
		break
	case "restore_app":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "RestoreApp requires 1 args")
			flag.Usage()
		}
		arg169 := flag.Arg(1)
		mbTrans170 := thrift.NewTMemoryBufferLen(len(arg169))
		defer mbTrans170.Close()
		_, err171 := mbTrans170.WriteString(arg169)
		if err171 != nil {
			Usage()
			return
		}
		factory172 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt173 := factory172.GetProtocol(mbTrans170)
		argvalue0 := admin.NewConfigurationRestoreRequest()
		err174 := argvalue0.Read(jsProt173)
		if err174 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.RestoreApp(context.Background(), value0))
		fmt.Print("\n")
		break
	case "start_partition_split":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "StartPartitionSplit requires 1 args")
			flag.Usage()
		}
		arg175 := flag.Arg(1)
		mbTrans176 := thrift.NewTMemoryBufferLen(len(arg175))
		defer mbTrans176.Close()
		_, err177 := mbTrans176.WriteString(arg175)
		if err177 != nil {
			Usage()
			return
		}
		factory178 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt179 := factory178.GetProtocol(mbTrans176)
		argvalue0 := admin.NewStartPartitionSplitRequest()
		err180 := argvalue0.Read(jsProt179)
		if err180 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.StartPartitionSplit(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_split_status":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QuerySplitStatus requires 1 args")
			flag.Usage()
		}
		arg181 := flag.Arg(1)
		mbTrans182 := thrift.NewTMemoryBufferLen(len(arg181))
		defer mbTrans182.Close()
		_, err183 := mbTrans182.WriteString(arg181)
		if err183 != nil {
			Usage()
			return
		}
		factory184 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt185 := factory184.GetProtocol(mbTrans182)
		argvalue0 := admin.NewQuerySplitRequest()
		err186 := argvalue0.Read(jsProt185)
		if err186 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QuerySplitStatus(context.Background(), value0))
		fmt.Print("\n")
		break
	case "control_partition_split":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ControlPartitionSplit requires 1 args")
			flag.Usage()
		}
		arg187 := flag.Arg(1)
		mbTrans188 := thrift.NewTMemoryBufferLen(len(arg187))
		defer mbTrans188.Close()
		_, err189 := mbTrans188.WriteString(arg187)
		if err189 != nil {
			Usage()
			return
		}
		factory190 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt191 := factory190.GetProtocol(mbTrans188)
		argvalue0 := admin.NewControlSplitRequest()
		err192 := argvalue0.Read(jsProt191)
		if err192 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ControlPartitionSplit(context.Background(), value0))
		fmt.Print("\n")
		break
	case "start_bulk_load":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "StartBulkLoad requires 1 args")
			flag.Usage()
		}
		arg193 := flag.Arg(1)
		mbTrans194 := thrift.NewTMemoryBufferLen(len(arg193))
		defer mbTrans194.Close()
		_, err195 := mbTrans194.WriteString(arg193)
		if err195 != nil {
			Usage()
			return
		}
		factory196 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt197 := factory196.GetProtocol(mbTrans194)
		argvalue0 := admin.NewStartBulkLoadRequest()
		err198 := argvalue0.Read(jsProt197)
		if err198 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.StartBulkLoad(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_bulk_load_status":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryBulkLoadStatus requires 1 args")
			flag.Usage()
		}
		arg199 := flag.Arg(1)
		mbTrans200 := thrift.NewTMemoryBufferLen(len(arg199))
		defer mbTrans200.Close()
		_, err201 := mbTrans200.WriteString(arg199)
		if err201 != nil {
			Usage()
			return
		}
		factory202 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt203 := factory202.GetProtocol(mbTrans200)
		argvalue0 := admin.NewQueryBulkLoadRequest()
		err204 := argvalue0.Read(jsProt203)
		if err204 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryBulkLoadStatus(context.Background(), value0))
		fmt.Print("\n")
		break
	case "control_bulk_load":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ControlBulkLoad requires 1 args")
			flag.Usage()
		}
		arg205 := flag.Arg(1)
		mbTrans206 := thrift.NewTMemoryBufferLen(len(arg205))
		defer mbTrans206.Close()
		_, err207 := mbTrans206.WriteString(arg205)
		if err207 != nil {
			Usage()
			return
		}
		factory208 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt209 := factory208.GetProtocol(mbTrans206)
		argvalue0 := admin.NewControlBulkLoadRequest()
		err210 := argvalue0.Read(jsProt209)
		if err210 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ControlBulkLoad(context.Background(), value0))
		fmt.Print("\n")
		break
	case "clear_bulk_load":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "ClearBulkLoad requires 1 args")
			flag.Usage()
		}
		arg211 := flag.Arg(1)
		mbTrans212 := thrift.NewTMemoryBufferLen(len(arg211))
		defer mbTrans212.Close()
		_, err213 := mbTrans212.WriteString(arg211)
		if err213 != nil {
			Usage()
			return
		}
		factory214 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt215 := factory214.GetProtocol(mbTrans212)
		argvalue0 := admin.NewClearBulkLoadStateRequest()
		err216 := argvalue0.Read(jsProt215)
		if err216 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.ClearBulkLoad(context.Background(), value0))
		fmt.Print("\n")
		break
	case "start_manual_compact":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "StartManualCompact requires 1 args")
			flag.Usage()
		}
		arg217 := flag.Arg(1)
		mbTrans218 := thrift.NewTMemoryBufferLen(len(arg217))
		defer mbTrans218.Close()
		_, err219 := mbTrans218.WriteString(arg217)
		if err219 != nil {
			Usage()
			return
		}
		factory220 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt221 := factory220.GetProtocol(mbTrans218)
		argvalue0 := admin.NewStartAppManualCompactRequest()
		err222 := argvalue0.Read(jsProt221)
		if err222 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.StartManualCompact(context.Background(), value0))
		fmt.Print("\n")
		break
	case "query_manual_compact":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "QueryManualCompact requires 1 args")
			flag.Usage()
		}
		arg223 := flag.Arg(1)
		mbTrans224 := thrift.NewTMemoryBufferLen(len(arg223))
		defer mbTrans224.Close()
		_, err225 := mbTrans224.WriteString(arg223)
		if err225 != nil {
			Usage()
			return
		}
		factory226 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt227 := factory226.GetProtocol(mbTrans224)
		argvalue0 := admin.NewQueryAppManualCompactRequest()
		err228 := argvalue0.Read(jsProt227)
		if err228 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.QueryManualCompact(context.Background(), value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}