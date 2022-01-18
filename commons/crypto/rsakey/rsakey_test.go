

package rsakey

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRsaPrivateKey1(t *testing.T) {

	priBytes, err := hex.DecodeString("308204a40201000282010100bd408cb5fdc87b712a031fc0ecd655078037084cccdef3315bf2720856961024447d58f455333df2b767ab099890f9d111f4d47931369cca4e756a6213c020a47d2324819837758f73c6e462dd7fd2424b4a7ba5f7d3d0b2266669a5c7ad4d3b9beb710fd5e06fb0f09091ae20350bb4aac62f5e7390ae77aab0718877a161730aaae7c2b21b7f9b79fb0358c13a5744fabb314d56d7de0dff08cbc378aa70ae3edacc804948f531563f7d94a10e6c30b4a722cf4b5fa96cd60b7b38faf140bb265ec73d0e4774d2fbb400c8bf456214b6bb7df763ca0f4041e64096bf3dcb0bafda8fb0fb00117b715ebce9691f85a03502dd6429009cd2487c55ccb0c76e3b020301000102820100562cd865de63b5e1f7b0f3e8936f9d54471a5a4b2e56af0260cdeb22e4da7c0b27acb41ebdd6e3bc6bcb26d4bdc5f61b3f43eea44d3422fcf3f0ff3a1da834b4f1ce58c7321abecb4d7ad97033500adbe910c770f1825fdb5f24ef226fe3f7f116b484cd2324897756735e029de3c9aac0c071bd5e7e1913a083ab8eea7a6fb67c4538a7cfb02a527549b1ef847ec28e855bf61a3b6488c52b1dd4d2025ea72ee2f576a227ab1b05e258e6a9e22a7e845d3f98e3b8b737990a507c0aac50c87aa5c479f52d4e3bd9094bf376e595a46792e6d2098a30a39ab1ec224b1c69f4ce1d4f3f28f2474d062331abb6fbb719c1fcbc1889f52649f33516244d9c7caa7102818100daa22fddbc7852f736273045fc914254507b248a063484db14b4fcdc00aa7942cfd953aa192d3cdbc75ab0267fa2a3569c7e01e2a410bde5151caa302bcb84c8d4f5f802a75cc69d6259f554ce42423856958ac6475b400b316d0b1762ad2c5fc31ceed5a5e95bfb5a0cd089bbcfb42e9309603f120ddf2ec4c58426c682954502818100dd98d8b95c4e24e19f2ff27e9d13ab5ca6ef4bdd65e399a04673b88864364e7edaa9c06b47a770e244b335c15e8c5c233c9ea2cab4037a57da1b37436154768d1fadbfafd4be70966c68b51c88d8e6aee5f3d7e8b93f1d5899c1c3a1fedc7a4c4ec3a9e13474341fc0a4f49ee36ca11a6b588b85a289213d55e4af560f3d6d7f02818100833755ed09a15981df3173ea7d241d20075170e399c7c978c71bbcaab98796d17f775a9c3b1208758b57256365b511bcf89d33ba776748e10563b7ccc36c191c839bc026af95a1ea714db64d18a171a6e86845eaac86da901d30e9b83653e2cef28619dd85fee162a0701274a79087fe6fe6efa9cac7228caa323517248ad8d5028180698f2cf6279d65fe406983b782b5e2f490e4ff1ba934a172f2fc9f1401c0c8e5aede1c363e7ce9ce2f71bb12b12a659db77bce0a8773fcacaace3a2613d03b65008930fdde771584e281827ce44786a41c106b72860425c39602f26151d9cf3c586ce698cbf6eaf9913842fb09552eed39e3851b1491044f8682187003747c9b028181009f0a1f98f9c291e1efb78ce8fd6f234ba65f37876fa43463811cf81fbe433da7e18483a736ee93d0fc2d08f2429db4f0a4363c2e2f0050c82dc9a26634b95f006b43d4f2feb4e7e8ef7b691d37f275e8b9f9f5e29e37374d502ba6a4648437d5e66dc0ced086eb7f78c75a45103ac38d23542e2051e8f2549d7f01f0dc130ffa")
	assert.Nil(t, err)

	rsaKey, err := NewRsaPrivateKey(priBytes)
	assert.Nil(t, err)
	assert.NotNil(t, rsaKey.PublicKey)
	assert.NotNil(t, rsaKey.PublicKey.RsaPublicKeyBytes)
	assert.NotNil(t, rsaKey.PublicKey.RsaPublicKey)
	assert.NotNil(t, rsaKey.PrivateKey)
	assert.NotNil(t, rsaKey.PrivateKeyBytes)

	assert.Equal(t, priBytes, rsaKey.PrivateKeyBytes)
}

func TestNewEmptyRsaPrivateKey(t *testing.T) {
	var priBytes []byte
	fmt.Println(len(priBytes))

	key, err := NewRsaPrivateKey(priBytes)
	assert.Nil(t, err)
	assert.Nil(t, key)
}

func TestNewErrorRsaPrivateKey(t *testing.T) {
	priBytes := []byte("aaaaaaaa")
	key, err := NewRsaPrivateKey(priBytes)
	assert.NotNil(t, err)
	assert.Nil(t, key)
}

func TestCreateRsaKey(t *testing.T) {
	rsaKey, err := CreateRsaKey()
	assert.Nil(t, err)
	encrypted, err := rsaKey.RsaEncrypt([]byte("test"))
	origData, err := rsaKey.RsaDecrypt(encrypted)
	assert.Equal(t, string(origData), "test")
}

func TestNewRsaPublicKey(t *testing.T) {
	rsaKey, err := CreateRsaKey()
	assert.Nil(t, err)
	rsaPublicKey, err := NewRsaPublicKey(rsaKey.PublicKey.RsaPublicKeyBytes)
	assert.Nil(t, err)
	assert.NotNil(t, rsaPublicKey)
}

func TestNewErrorRsaPublicKey(t *testing.T) {
	rsaPublicKey, err := NewRsaPublicKey([]byte("test"))
	assert.Nil(t, rsaPublicKey)
	assert.NotNil(t, err)
}

func TestNewRsaKey(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "testnewRsaKey")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	rsaKey, err := NewRsaKey(tmpdir)
	assert.Nil(t, err)
	assert.NotNil(t, rsaKey)
}
