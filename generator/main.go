/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"unicode"
)

func privateMethodName(methodName string) string {
	return string(unicode.ToLower(rune(methodName[0]))) + methodName[1:]
}

func generateAdminRPC(rpcCode, rpcName, reqName, respName string) string {
	s := "\n"
	s += fmt.Sprintf("func (ms *metaSession) %s(ctx context.Context, req *admin.%s) (*admin.%s, error) {\n", privateMethodName(rpcName), reqName, respName)
	s += fmt.Sprintf("	arg := admin.NewAdminClient%sArgs()\n", rpcName)
	s += fmt.Sprintf("	arg.Req = req\n")
	s += fmt.Sprintf("	result, err := ms.call(ctx, arg, \"%s\")\n", rpcCode)
	s += fmt.Sprintf("	if err != nil {\n")
	s += fmt.Sprintf("		return nil, fmt.Errorf(\"RPC to session %%s failed: %%s\", ms, err)\n")
	s += fmt.Sprintf("	}\n")
	s += fmt.Sprintf("	ret, _ := result.(*admin.AdminClient%sResult)\n", rpcName)
	s += fmt.Sprintf("	return ret.GetSuccess(), nil\n")
	s += fmt.Sprintf("}\n")
	s += "\n"
	s += fmt.Sprintf("// %s is auto-generated\n", rpcName)
	s += fmt.Sprintf("func (m *MetaManager) %s(ctx context.Context, req *admin.%s) (*admin.%s, error) {\n", rpcName, reqName, respName)
	s += fmt.Sprintf("	resp, err := m.call(ctx, func(rpcCtx context.Context, ms *metaSession) (metaResponse, error) {\n")
	s += fmt.Sprintf("		return ms.%s(rpcCtx, req)\n", privateMethodName(rpcName))
	s += fmt.Sprintf("	})\n")
	s += fmt.Sprintf("	if err == nil {\n")
	s += fmt.Sprintf("		if resp.GetErr().Errno != base.ERR_OK.String() {\n")
	s += fmt.Sprintf("			return resp.(*admin.%s), fmt.Errorf(\"%s failed: %%s\", resp.GetErr().String())\n", respName, rpcName)
	s += fmt.Sprintf("		}\n")
	s += fmt.Sprintf("		return resp.(*admin.%s), nil\n", respName)
	s += fmt.Sprintf("	}\n")
	s += fmt.Sprintf("	return nil, err\n")
	s += fmt.Sprintf("}\n")
	return s
}

func generateHeader() {
	fmt.Print(`/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Code generated by "generator -i=admin.csv > admin_rpc_types.go"; DO NOT EDIT.
package session

`)
}

func generateAllAdminRPC() {
	content, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input file \"%s\" : %s\n", inputFilePath, err)
		flag.Usage()
		os.Exit(1)
	}

	csvReader := csv.NewReader(bytes.NewReader(content))
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse csv file \"%s\" : %s\n", inputFilePath, err)
		os.Exit(1)
	}

	generateHeader()
	fmt.Print(`
import (
	"context"
	"fmt"

	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/admin"
	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/base"
)
`)

	for i, lineArgs := range lines {
		if len(lineArgs) != 4 {
			fmt.Fprintf(os.Stderr, "invalid arguments in line %d\n", i)
			os.Exit(1)
		}
		s := generateAdminRPC(lineArgs[0], lineArgs[1], lineArgs[2], lineArgs[3])
		fmt.Print(s)
	}
}

func generateRAdminRPC(rpcCode, rpcName, reqName, respName string) string {
	s := "\n"
	s += fmt.Sprintf("// %s is auto-generated\n", rpcName)
	s += fmt.Sprintf("func (rs *ReplicaSession) %s(ctx context.Context, req *radmin.%s) (*radmin.%s, error) {\n", rpcName, reqName, respName)
	s += fmt.Sprintf("	arg := radmin.NewReplicaClient%sArgs()\n", rpcName)
	s += fmt.Sprintf("	arg.Req = req\n")
	s += fmt.Sprintf("	result, err := rs.CallWithGpid(ctx, &base.Gpid{Appid: 0, PartitionIndex: 0}, 0, arg, \"%s\")\n", rpcCode)
	s += fmt.Sprintf("	if err == nil {\n")
	s += fmt.Sprintf("		ret, _ := result.(*radmin.ReplicaClient%sResult)\n", rpcName)
	s += fmt.Sprintf("		resp := ret.GetSuccess()\n")
	s += fmt.Sprintf("		if resp.GetErr().Errno != base.ERR_OK.String() {\n")
	s += fmt.Sprintf("			return resp, fmt.Errorf(\"%s to session %%s failed: %%s\", rs, resp.GetErr().String())\n", rpcName)
	s += fmt.Sprintf("		}\n")
	s += fmt.Sprintf("		return resp, nil\n")
	s += fmt.Sprintf("	}\n")
	s += fmt.Sprintf("	return nil, err\n")
	s += fmt.Sprintf("}\n")
	return s
}

func generateAllRAdminRPC() {
	content, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input file \"%s\" : %s\n", inputFilePath, err)
		flag.Usage()
		os.Exit(1)
	}

	csvReader := csv.NewReader(bytes.NewReader(content))
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse csv file \"%s\" : %s\n", inputFilePath, err)
		os.Exit(1)
	}

	generateHeader()
	fmt.Print(`
import (
	"context"
	"fmt"

	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/base"
	"github.com/nanorth/pegasus-go-client-for-meta-proxy/idl/radmin"
)
`)

	for i, lineArgs := range lines {
		if len(lineArgs) != 4 {
			fmt.Fprintf(os.Stderr, "invalid arguments in line %d\n", i)
			os.Exit(1)
		}
		s := generateRAdminRPC(lineArgs[0], lineArgs[1], lineArgs[2], lineArgs[3])
		fmt.Print(s)
	}
}

var inputFilePath string
var rpcSuite string

func init() {
	flag.StringVar(&inputFilePath, "i", "", "the RPC definition csv file path")

	flag.StringVar(&rpcSuite, "t", "", "the RPC suite to generate, options [\"admin\", \"radmin\"]")
}

func main() {
	flag.Parse()

	if rpcSuite == "admin" {
		generateAllAdminRPC()
	} else if rpcSuite == "radmin" {
		generateAllRAdminRPC()
	}
}
