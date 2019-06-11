// Putting some sample data.

package service

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	pb "github.com/maelvls/quote/schema/user"
)

// LoadSampleUsers loads some hard-coded users into database.
func (svc *UserImpl) LoadSampleUsers() error {
	var users []pb.User
	err := json.Unmarshal(sampleUsers, &users)
	if err != nil {
		return fmt.Errorf("could not parse json: %v", err)
	}

	txn := svc.DB.Txn(true)
	for _, user := range users {
		u := user
		if err := txn.Insert("user", &u); err != nil {
			return err
		}
	}
	txn.Commit()
	logrus.Debugf("added user samples to DB")

	return nil
}

// SampleUsers is a json file of users for testing purposes.
var sampleUsers = []byte(`
[
  {
    "id": "5cfdf218f7efd273906c5b9e",
    "age": 51,
    "name": {
      "first": "Valencia",
      "last": "Dorsey"
    },
    "email": "valencia.dorsey@email.info",
    "phone": "+1 (906) 568-2594",
    "address": "941 Merit Court, Grill, Mississippi, 4961"
  },
  {
    "id": "5cfdf218090eae728f3ebf2d",
    "age": 52,
    "name": {
      "first": "Brianna",
      "last": "Shelton"
    },
    "email": "brianna.shelton@email.org",
    "phone": "+1 (814) 482-3880",
    "address": "255 Cortelyou Road, Volta, Indiana, 1608"
  },
  {
    "id": "5cfdf2186e34d5ec8e605988",
    "age": 60,
    "name": {
      "first": "Snider",
      "last": "Fisher"
    },
    "email": "snider.fisher@email.biz",
    "phone": "+1 (918) 591-2784",
    "address": "363 Williamsburg Street, Chicopee, Illinois, 6316"
  },
  {
    "id": "5cfdf21883e020967de82837",
    "age": 48,
    "name": {
      "first": "Pacheco",
      "last": "Fitzgerald"
    },
    "email": "pacheco.fitzgerald@email.name",
    "phone": "+1 (828) 442-3262",
    "address": "278 McKibben Street, Nicholson, South Dakota, 3793"
  },
  {
    "id": "5cfdf218862e0be14633412a",
    "age": 35,
    "name": {
      "first": "Brock",
      "last": "Stanley"
    },
    "email": "brock.stanley@email.me",
    "phone": "+1 (836) 594-3347",
    "address": "748 Aster Court, Elwood, Guam, 7446"
  },
  {
    "id": "5cfdf2185976a696e86279a1",
    "age": 42,
    "name": {
      "first": "Hardin",
      "last": "Patton"
    },
    "email": "hardin.patton@email.com",
    "phone": "+1 (977) 536-2989",
    "address": "241 Russell Street, Robinson, Oregon, 9576"
  },
  {
    "id": "5cfdf21851279432185e9811",
    "age": 26,
    "name": {
      "first": "Walter",
      "last": "Prince"
    },
    "email": "walter.prince@email.co.uk",
    "phone": "+1 (804) 553-3262",
    "address": "204 Ralph Avenue, Gibbsville, Michigan, 6698"
  },
  {
    "id": "5cfdf218d7b6ae5366adfb8e",
    "age": 22,
    "name": {
      "first": "Acevedo",
      "last": "Quinn"
    },
    "email": "acevedo.quinn@email.us",
    "phone": "+1 (886) 442-2144",
    "address": "403 Lawn Court, Walland, Federated States Of Micronesia, 8260"
  },
  {
    "id": "5cfdf2180a1645ce50cd9aa2",
    "age": 28,
    "name": {
      "first": "Billie",
      "last": "Norton"
    },
    "email": "billie.norton@email.io",
    "phone": "+1 (934) 524-3718",
    "address": "699 Rapelye Street, Dupuyer, Ohio, 4175"
  },
  {
    "id": "5cfdf218a9a0c61919af13a6",
    "age": 51,
    "name": {
      "first": "Solis",
      "last": "Irwin"
    },
    "email": "solis.irwin@email.tv",
    "phone": "+1 (855) 413-3330",
    "address": "739 Poly Place, Rosburg, Colorado, 9608"
  },
  {
    "id": "5cfdf2188df89b48900fd70f",
    "age": 48,
    "name": {
      "first": "Wilkerson",
      "last": "Mosley"
    },
    "email": "wilkerson.mosley@email.biz",
    "phone": "+1 (884) 464-2806",
    "address": "734 Kosciusko Street, Marbury, Connecticut, 3037"
  },
  {
    "id": "5cfdf21804725eefa8d9ec69",
    "age": 33,
    "name": {
      "first": "Alford",
      "last": "Cole"
    },
    "email": "alford.cole@email.net",
    "phone": "+1 (822) 589-2083",
    "address": "763 Halleck Street, Elbert, Nevada, 3291"
  },
  {
    "id": "5cfdf2184ad9cd19459891f7",
    "age": 31,
    "name": {
      "first": "Stone",
      "last": "Briggs"
    },
    "email": "stone.briggs@email.info",
    "phone": "+1 (828) 438-2266",
    "address": "531 Atkins Avenue, Neahkahnie, Tennessee, 3981"
  },
  {
    "id": "5cfdf21826f2f78ece771e03",
    "age": 57,
    "name": {
      "first": "Ratliff",
      "last": "Herring"
    },
    "email": "ratliff.herring@email.org",
    "phone": "+1 (949) 540-2608",
    "address": "246 Greene Avenue, Blairstown, Puerto Rico, 6855"
  },
  {
    "id": "5cfdf218a689729c23f25847",
    "age": 48,
    "name": {
      "first": "Angeline",
      "last": "Stokes"
    },
    "email": "angeline.stokes@email.biz",
    "phone": "+1 (970) 569-3963",
    "address": "526 Java Street, Hailesboro, Pennsylvania, 1648"
  },
  {
    "id": "5cfdf218d025ace57aadcc01",
    "age": 56,
    "name": {
      "first": "Santos",
      "last": "Slater"
    },
    "email": "santos.slater@email.name",
    "phone": "+1 (858) 533-2802",
    "address": "459 Sharon Street, Belleview, Kentucky, 5483"
  },
  {
    "id": "5cfdf218b2a4b08ad8efd775",
    "age": 35,
    "name": {
      "first": "Ina",
      "last": "Perkins"
    },
    "email": "ina.perkins@email.me",
    "phone": "+1 (844) 507-2552",
    "address": "899 Miami Court, Temperanceville, Virginia, 2821"
  },
  {
    "id": "5cfdf218e5f9edbd4bba3faf",
    "age": 46,
    "name": {
      "first": "Rice",
      "last": "Pierce"
    },
    "email": "rice.pierce@email.com",
    "phone": "+1 (899) 428-2988",
    "address": "291 Boardwalk , Chloride, North Carolina, 8401"
  },
  {
    "id": "5cfdf218d09827cb4530b5d7",
    "age": 58,
    "name": {
      "first": "Shields",
      "last": "Moody"
    },
    "email": "shields.moody@email.co.uk",
    "phone": "+1 (953) 554-3038",
    "address": "350 Powell Street, Chaparrito, Massachusetts, 2556"
  },
  {
    "id": "5cfdf218387d3edf16da6a46",
    "age": 52,
    "name": {
      "first": "Jenifer",
      "last": "Valencia"
    },
    "email": "jenifer.valencia@email.us",
    "phone": "+1 (988) 463-2789",
    "address": "948 Jefferson Street, Guthrie, Louisiana, 2483"
  },
  {
    "id": "5cfdf218e8491ba6b28c17bf",
    "age": 56,
    "name": {
      "first": "Beasley",
      "last": "Byrd"
    },
    "email": "beasley.byrd@email.io",
    "phone": "+1 (819) 597-2912",
    "address": "213 McKibbin Street, Veguita, New Jersey, 3943"
  },
  {
    "id": "5cfdf2182b8e2574925daa7c",
    "age": 21,
    "name": {
      "first": "Helen",
      "last": "Walker"
    },
    "email": "helen.walker@email.tv",
    "phone": "+1 (805) 518-2099",
    "address": "861 Conselyea Street, Elliott, Texas, 4229"
  },
  {
    "id": "5cfdf218b0ae76504da8a23c",
    "age": 22,
    "name": {
      "first": "Ivy",
      "last": "Stephens"
    },
    "email": "ivy.stephens@email.biz",
    "phone": "+1 (948) 401-2314",
    "address": "246 Bushwick Avenue, Grazierville, California, 4664"
  },
  {
    "id": "5cfdf2183d741104c707c5d9",
    "age": 31,
    "name": {
      "first": "Benjamin",
      "last": "Frazier"
    },
    "email": "benjamin.frazier@email.net",
    "phone": "+1 (953) 407-3166",
    "address": "289 Cyrus Avenue, Templeton, Maine, 5964"
  },
  {
    "id": "5cfdf218165c21a887626054",
    "age": 59,
    "name": {
      "first": "Hodge",
      "last": "Cabrera"
    },
    "email": "hodge.cabrera@email.info",
    "phone": "+1 (923) 543-3169",
    "address": "521 Richards Street, Takilma, Missouri, 4287"
  },
  {
    "id": "5cfdf218ad60c3c443d2f579",
    "age": 51,
    "name": {
      "first": "Kent",
      "last": "Cochran"
    },
    "email": "kent.cochran@email.org",
    "phone": "+1 (945) 512-2231",
    "address": "803 Cranberry Street, Inkerman, Marshall Islands, 6929"
  },
  {
    "id": "5cfdf218b65a0f7cfb60fddc",
    "age": 44,
    "name": {
      "first": "Noreen",
      "last": "Parks"
    },
    "email": "noreen.parks@email.biz",
    "phone": "+1 (950) 461-3686",
    "address": "872 Milford Street, Goldfield, Minnesota, 3340"
  },
  {
    "id": "5cfdf2185aae57eaf6d8ffaa",
    "age": 24,
    "name": {
      "first": "Marion",
      "last": "Zimmerman"
    },
    "email": "marion.zimmerman@email.name",
    "phone": "+1 (903) 437-2904",
    "address": "731 Jamison Lane, Independence, North Dakota, 846"
  },
  {
    "id": "5cfdf2181cd125fe11379967",
    "age": 23,
    "name": {
      "first": "Blanca",
      "last": "Lang"
    },
    "email": "blanca.lang@email.me",
    "phone": "+1 (848) 458-2687",
    "address": "995 Meadow Street, Greenbackville, New Mexico, 1237"
  },
  {
    "id": "5cfdf2181fc77ecd79910093",
    "age": 27,
    "name": {
      "first": "Dawson",
      "last": "Boyer"
    },
    "email": "dawson.boyer@email.com",
    "phone": "+1 (804) 566-3741",
    "address": "283 Jewel Street, Salvo, Oklahoma, 1417"
  }
]
`)
