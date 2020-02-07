echo "Election Transaction:"

curl --header "Content-Type: application/json" --request POST --data '{
	"type": "Election",
	"electionName": "Student Government Elections",
	"institutionName": "West Chester University",
	"description": "Spring Elections",
	"startDate": "Thu, 06 Feb 2020 22:26:16 EST",
	"endDate": "Mon, 24 Aug 2020 22:26:16 EDT",
	"emailDomain": "wcupa.edu",
	"positions": [
		{
			"id": "demfrmeororev",
			"displayName": "Student Government President",
			"candidates": [
				{
					"name": "John Doe",
					"key": "test1"
				},
				{
					"name": "Sarah Jennings",
					"key": "test2"
				},
				{
					"name": "Maximus Footless",
					"key": "test3"
				}
			]
		}
	],
	"sender": "3059301306072a8648ce3d020106082a8648ce3d030107034200048c2be6467d4a477ac8b5cbbded6528af7b6c44291853467448e585a4e57e3c7cdcb52646d192959a54c770f2c79cb6e7a0c3b716275588b4e7433aeb0128eac2",
	"signature": "7b2252223a31313031393736373032313534373536393138353736333732363931373032393138353033353038343439373235313436333136343638303730343036323336303636383034373138333236342c2253223a36343330313332303239383833353831323037343331333139313437303834343233313331333135323331363438353134323930383132373838343438333835393332323035323134333438357d"
}' http://localhost:7003/election



echo "Registration Transaction:"

curl --header "Content-Type: application/json" --request POST --data '{
	"emailAddress": "jacob@neubaum.com",
	"firstName": "Jacob",
	"lastName": "Neubaum",
	"dateOfBirth": "3/9/1999",
	"electionName": "7b2252223a31313031393736373032313534373536393138353736333732363931373032393138353033353038343439373235313436333136343638303730343036323336303636383034373138333236342c2253223a36343330313332303239383833353831323037343331333139313437303834343233313331333135323331363438353134323930383132373838343438333835393332323035323134333438357d",
	"electionAdmin": "3059301306072a8648ce3d020106082a8648ce3d030107034200048c2be6467d4a477ac8b5cbbded6528af7b6c44291853467448e585a4e57e3c7cdcb52646d192959a54c770f2c79cb6e7a0c3b716275588b4e7433aeb0128eac2",
	"publicKey": "3059301306072a8648ce3d020106082a8648ce3d0301070342000414aa5570c4d61f16dbb0d19fb8759b1630c6952cb364dd19d19cbf12db8d62cf44e6b5431e1fb7e24d5f0de52e64410a7509cc93b02c760a2e5e4e7e4d267315",
	"senderSig": "",
	"code": "",
	"timestamp": ""
}' http://localhost:7003/election/test/register



echo "Vote Transaction:"

curl --header "Content-Type: application/json" --request POST --data '{
	"type": "Vote",
	"electionName": "7b2252223a31313031393736373032313534373536393138353736333732363931373032393138353033353038343439373235313436333136343638303730343036323336303636383034373138333236342c2253223a36343330313332303239383833353831323037343331333139313437303834343233313331333135323331363438353134323930383132373838343438333835393332323035323134333438357d",
	"receivers": [
		{
			"id": "demfrmeororev",
			"key": "test1"
		}
	],
	"sender": "3059301306072a8648ce3d020106082a8648ce3d0301070342000414aa5570c4d61f16dbb0d19fb8759b1630c6952cb364dd19d19cbf12db8d62cf44e6b5431e1fb7e24d5f0de52e64410a7509cc93b02c760a2e5e4e7e4d267315",
	"signature": "7b2252223a36323032383733333738363636303530333730323037363436393634393237393035353036363132323337323837363232333435373636323938303839333531333538343630383134393838342c2253223a39393539383233343233353731333833343738343730313630343232393031373332383334333932373233363630383530333639373739363130383038303133353939313234333337363631357d"
}' http://localhost:7003/election/test/vote



