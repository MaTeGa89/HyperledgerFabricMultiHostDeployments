const axios = require("axios");
const url = "http://34.123.162.252:4000/channels/mychannel/chaincodes/supply-chain";

const login = async () => {
	let d = {
		"username": "akshay",
		"orgName": "Org1"
	}
	try {
		let resp = await axios.post("http://34.123.162.252:4000/users", d, { headers: { "Content-Type": "application/json" } })
		return resp.data.token
	} catch (error) {
		return error
	}
}


const test = async () => {
	let token = await login()
	console.log(token)
}

// test()

const getRandomTemperature = (min, max) => {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}



const updateInvoiceBatch = async (batchId, sensorId) => {

	let token = await login()
	console.log(token)

	let temperatureLocationData

	setInterval(() => {
		console.log("inside")
			 temperatureLocationData = {
				temperature : getRandomTemperature(0,27).toString(),
				timestamp: (+ new Date()).toString(),
				temperatureSensorId: sensorId,
				longitude: (Math.random()*360 - 180).toString(),
				latitude: (Math.random()*360 - 180).toString()

			}
			args = [
				batchId,
				JSON.stringify(temperatureLocationData)
			]

			let data = {
				"fcn": "UpdateVaccineBatch",
				"chaincodeName": "supply-chain",
				"channelName": "mychannel",
				"args": args
			}


			return axios.post(url, data, {
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json"
				}
			}
			).then(resp => {
				console.log(resp.data.result)
			}
			).catch(function (error) { console.log(error); });


	}, 3000)

};

updateInvoiceBatch("5", "1234")
