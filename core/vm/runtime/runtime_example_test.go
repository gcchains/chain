

package runtime_test

import (
	"fmt"

	"github.com/gcchains/chain/core/vm/runtime"
	"github.com/ethereum/go-ethereum/common"
)

func ExampleExecute() {
	ret, _, err := runtime.Execute(common.Hex2Bytes("6060604052600a8060106000396000f360606040526008565b00"), nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
	// Output:
	// [96 96 96 64 82 96 8 86 91 0]
}
