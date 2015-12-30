package yapl
import "testing"

func TestYaplExec(t *testing.T) {

}

var moduleA = `
module a {
	prefix "";
	namespace "";
	revision 0;
	leaf f {
		type string;
	}
	list b {
		key "c";
		leaf c {
			type string;
		}
		container d {
			leaf e {
				type string;
			}
		}
	}
}
`

var moduleZ = `
module z {
	prefix "";
	namespace "";
	revision 0;
	leaf u {
		type string;
	}
	list y {
		key "x";
		leaf x {
			type string;
		}
		container w {
			leaf v {
				type string;
			}
		}
	}
}
`