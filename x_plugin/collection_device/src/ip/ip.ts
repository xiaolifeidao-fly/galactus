import { getData } from "@src/utils/axios";



export async function getIp() {
    return await getData(String, "/ip/collectDeviceIp")
}
