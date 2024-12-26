const jsdom = require("jsdom");
const {JSDOM} = jsdom;
const dom = new JSDOM(`<!DOCTYPE html><p>Hello world</p>`, {
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
});



function get_enviroment(proxy_array) {
    for (var i = 0; i < proxy_array.length; i++) {
        handler = '{\n' +
            '    get: function(target, property, receiver) {\n' +
            '        console.log("　　　　:", "get  ", "　　　　:", ' +
            '"' + proxy_array[i] + '" ,' +
            '"  　　　　:", property, ' +
            '"  　　　　　　　　:", ' + 'typeof property, ' +
            // '"  　　　　　:", ' + 'target[property], ' +
            '"  　　　　　　　　　:", typeof target[property]);\n' +
            '        return target[property];\n' +
            '    },\n' +
            '    set: function(target, property, value, receiver) {\n' +
            '        console.log("　　　　:", "set  ", "　　　　:", ' +
            '"' + proxy_array[i] + '" ,' +
            '"  　　　　:", property, ' +
            '"  　　　　　　　　:", ' + 'typeof property, ' +
            // '"  　　　　　:", ' + 'target[property], ' +
            '"  　　　　　　　　　:", typeof target[property]);\n' +
            '        return Reflect.set(...arguments);\n' +
            '    }\n' +
            '}'
        eval('try{\n' + proxy_array[i] + ';\n'
            + proxy_array[i] + '=new Proxy(' + proxy_array[i] + ', ' + handler + ')}catch (e) {\n' + proxy_array[i] + '={};\n'
            + proxy_array[i] + '=new Proxy(' + proxy_array[i] + ', ' + handler + ')}')
    }
}

proxy_array = ['window', 'document', 'location', 'navigator', 'history', 'screen']  // 　　　　　　　　　　
window = global;
// document = dom.window.document
window.requestAnimationFrame = function (res) {
    console.log("window.requestAnimationFrame ->", res)
}
window._sdkGlueVersionMap = {
    "sdkGlueVersion": "1.0.1.17.01",
    "bdmsVersion": "1.0.1.17.01",
    "captchaVersion": "1.0.0.63.01"
}
window.onwheelx = {
    "_Ax": "0X21"
}
window.innerWidth = 1707
window.innerHeight = 791
window.outerWidth = 1707
window.outerHeight = 912
window.screenX = 0
window.screenY = 0
window.pageYOffset = 0
window.fetch = function (res) {
    console.log("window.fetch ->",res)
}

screen = {
    availWidth: 1707,
    availHeight: 912,
    width: 1707,
    height: 960,
    colorDepth: 24,
    pixelDepth: 24,
}
XMLHttpRequest = dom.window.XMLHttpRequest;
span = []
fpk1="U2FsdGVkX1+9wiJGncRAWHzVOFI4HswOwm2Z3Q1xgYMYHW62lWpvk0WbOOSp2yNrrtGNPy84Da7TIo0GFy2O+Q==";
byte_fid="23439c627382f6a0e13b1f54e7a81c7f";
UIFID="96cd3b166f3029d7c1cc3f64582454ab8a83ff1f9e6d6689076dd47ef1dca5f838478297937fcb88316a6ee9493aba5f8d54cd67ea1d1784b59d9556552d13f105e3f5242002dfd9e2235682f0ed8bc2bca11f5182cadab8b57e511e045ec0bb4de6fe3074b07c589e096bbf2cb03bc38cceaa7ff28af647ba70fe054ab39588c42e7c7050c93bdcdfd2dc7daffe47f32717bda207d1a5d22c8e4bc7c1efafc9";


setTimeout = function () {
}
navigator = {
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
    vendorSubs: {
        "ink": 1713882586968
    },
    platform: 'Win32'
}

get_enviroment(proxy_array)
console.log("ddd")

