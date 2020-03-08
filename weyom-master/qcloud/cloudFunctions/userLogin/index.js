/* eslint-disable import/newline-after-import */
/* eslint-disable import/no-commonjs */
// 云函数入口文件
const cloud = require('wx-server-sdk')
const rp = require('request-promise')
cloud.init({
  env: ''
})
// 云函数入口函数
exports.main = async event => {
    console.log(event.code)
    let { OPENID, APPID, UNIONID } = cloud.getWXContext()
    console.log(OPENID, APPID, UNIONID)
    return {
     OPENID,
     APPID,
     UNIONID,
    }
  }
