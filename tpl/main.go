/*
Copyright © 2024 Spring Zhang spring.zhang@zettablock.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tpl

func ConfigTemplate() []byte {
	return []byte(`specVersion: 0.0.1 # not used for now
# required, usually this should be company name without space.
# if metadata.schema not provided, this will be used as metadata db schema
org: zb # no space and specifia char allowed
kind: ethereum # chain name
network: holesky # chain network
version: 0 # not used for now
name: demo_testnet_holesky_project1 # required, not space or specical chars allowed
source:
	schema: "ethereum_holesky"
	sourceDB: SELF_HOSTED_PROD_DB # in staging and prod, we can use prod db with read-only user
	startBlock: 1167044
#  addresses:
#    - "0x055733000064333CaDDbC92763c58BF0192fFeBf" # testnet holesky AVS-Directory contract
metadata:
	metadataDB: RDS_TEST_DB
	schema: zb
destination:
	destinationDB: RDS_TEST_DB
	schema: "zb" # schema name for destination
	eventHandlers: # multiple handlers
	- event: Transfer
	handler: HandleTransfer
blockHandlers:
	- handler: HandleBlock
`)
}

func BlockHandlersTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
	"time"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
)

// *utils.Deps contains *gorm.DB which can be used for CRUD
// *utils.Deps also contains Logger
func HandleBlock(block ethereum.Block, deps *utils.Deps) (bool, error) {
	deps.Logger.Info("HandleBlock", "block number", block.Number, "pipeline_name", deps.Config.GetPipelineName())
	fmt.Println("HandleBlock called in proj3")
	fmt.Println("block", block)
	fmt.Println("block number", block.Number)
	fmt.Println("block hash", block.Hash)
	fmt.Println("block timestamp", block.Timestamp)
	fmt.Println("block date", block.BlockDate)
	time.Sleep(10 * time.Second)
	return false, nil
}
`)
}

func EventHandlersTemplate() []byte {
	return []byte(`package main

import (
	"fmt"
	"time"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
)

func HandleTransfer(log ethereum.Log, deps *utils.Deps) (bool, error) {
	deps.Logger.Info("HandleTransfer", "block number", log.BlockNumber, "pipeline_name", deps.Config.GetPipelineName())
	fmt.Println("HandleTransfer called in proj3")
	fmt.Println("log", log)
	fmt.Println("block number", log.BlockNumber)
	fmt.Println("log index", log.LogIndex)
	fmt.Println("log block time", log.BlockTime)
	fmt.Println("log block date", log.BlockDate)
	fmt.Println("log event", log.Event)
	time.Sleep(2 * time.Second)
	return false, nil
}	
`)
}
