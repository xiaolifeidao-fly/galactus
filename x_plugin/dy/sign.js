


class AcSign{

    constructor() {
        this.finalNum = 0;
    }

    bigCountOperation(str) {
        let charCodeCount;
        for (let i = 0; i < str.length; i++) {
            charCodeCount = str.charCodeAt(i)
            this.finalNum = ((this.finalNum ^ charCodeCount) * 65599) >>> 0
        }
        return this.finalNum
    }
    
    countToText(deciNum, ac_signature) {
        let offList = [24, 18, 12, 6, 0]
        offList.forEach(function (value) {
            let keyNum = deciNum >> value & 63
            let valNum = keyNum < 26 ? 65 : keyNum < 52 ? 71 : keyNum < 62 ? -4 : -17
            let asciiCode = keyNum + valNum
            ac_signature += String.fromCharCode(asciiCode)
        })
        return ac_signature
    }

    loadAcSignature(url, ac_nonce, ua) {
        var temp = 0;
        let timeStamp = new Date().getTime().toString()
        this.bigCountOperation(timeStamp)
        let urlNum = this.bigCountOperation(url)
        let longStr = ((65521 * (this.finalNum % 65521) ^ timeStamp) >>> 0).toString(2)
        while (longStr.length !== 32) {
            longStr = "0" + longStr
        }
        let binaryNum = "10000000110000" + longStr
        let deciNum = parseInt(binaryNum, 2)
        var ac_signature = "_02B4Z6wo00f01";
        ac_signature = this.countToText(deciNum >> 2, ac_signature)
        ac_signature = this.countToText(deciNum << 28 | 515, ac_signature)
        ac_signature = this.countToText((deciNum ^ 1489154074) >>> 6, ac_signature)
        let aloneNum = (deciNum ^ 1489154074) & 63
        let aloneVal = aloneNum < 26 ? 65 : aloneNum < 52 ? 71 : aloneNum < 62 ? -4 : -17
        ac_signature += String.fromCharCode(aloneNum + aloneVal)
    
        this.finalNum = 0
        let deciOperaNum = this.bigCountOperation(deciNum.toString())
        let nonceNum = this.bigCountOperation(ac_nonce)
        this.finalNum = deciOperaNum
        this.bigCountOperation(ua)
        ac_signature = this.countToText((nonceNum % 65521 | (this.finalNum % 65521 << 16)) >> 2, ac_signature)
    
        ac_signature = this.countToText((((this.finalNum % 65521 << 16) ^ (nonceNum % 65521)) << 28) | ((deciNum << 524576 ^ 524576) >>> 4), ac_signature)
    
        ac_signature = this.countToText(urlNum % 65521, ac_signature)
    
        let _str = ac_signature;
    
        for (let i of _str) {
            temp = ((temp * 65599) + i.charCodeAt(0)) >>> 0
        }
        let lastStr = temp.toString(16)
        ac_signature += lastStr.slice(lastStr.length - 2, lastStr.length)
        return ac_signature
    }
}
function getAcSign(url, ac_nonce, ua){
    return new AcSign().loadAcSignature(url, ac_nonce, ua);
}

module.exports = {
    getAcSign
};