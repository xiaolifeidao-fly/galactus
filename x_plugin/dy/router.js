const {use} = require("./application")

const aBogus = require('./a_bogus_192')
const sign = require('./sign')



function getAcSign(req, res){
    const body = req.body;
    const url = body['url'];
    const acNonce = body['acNonce'];
    const ua = body['ua'];
    const acSignature = sign.getAcSign(url, acNonce, ua);
    res.send({"acSignature": acSignature});
}

function getAbogus(req, res){
    const body = req.body;
    const params = body['params'];
    // console.log("req.body    ", body);
    const ua = body['ua'];
    var data = aBogus.getAbogus(params, ua);
    res.send({"aBogus":data});
}


function init(){
    use('/dy/ac/sign', getAcSign);
    use('/dy/abogus/sign', getAbogus);

}

module.exports = {
    init
}