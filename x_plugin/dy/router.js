const {use} = require("../application")

const aBogus = require('./a_bogus')
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
    const ua = body['ua'];
    var data = aBogus.generate_a_bogus(params, ua);
    res.send({"aBogus":data});
}


function init(){
    use('/dy/ac/sign', getAcSign);
    use('/dy/abogus/sign', getAbogus);

}

module.exports = {
    init
}