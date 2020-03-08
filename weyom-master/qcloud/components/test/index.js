import Taro, {Component} from '@tarojs/taro'
import {View, Text, Image, Checkbox} from '@tarojs/components'
import {AtButton, AtToast, AtInput, AtMessage, AtCard} from 'taro-ui'
import './index.scss'
import usericon from "../../../static/icon/mine-selected.png";
import cardicon from "../../../static/icon/home-selected.png";
import {StorageKey, URL_ESIGN} from "../../../constants";
import {
  getOpenId,
  requestPostm,
  showMsg,
  getUserInfoDB,
} from "../../../utils_other";
import {protocol} from "./content";
// import erroricon from "../../static/image/error.png";
import erroricon from "../../static/image/auth_fail.png";

/**
 * 失败页面
 */
export default class realname_fail extends Component {
  config = {
    navigationBarTitleText: "借条服务协议",
    // navigationStyle: 'custom'
  }

  constructor() {
    super(...arguments)

    this.state = {

      errTitle: '出错了', // 错误文字

      code: '',
      message: '',

      // 消息通知组件
      msgText: '',
      msgType: 'info',
      showMsgView: false,

    }
  }


  componentDidShow() {
  }

  componentWillMount() {
  }

  componentDidMount() {
    this.checkAuth()
  }


  componentWillUnmount() {

  }

  componentDidHide() {
    this.cleanData()
  }

  checkAuth() {
    Taro.getSetting({
      success(res) {
        console.log(res.authSetting)
        if (!res.authSetting['scope.userInfo']) {
          Taro.navigateTo({
            url: '/subpages/auth/authUser/index'
          })
        }
      }
    })
  }

  goHomePage() {
    Taro.reLaunch({
      url: '/pages/index/index'
    })
  }

  realNameAgain() {
    Taro.navigateBack({
      delta: 1
    })
  }


  render() {
    return (
      <View className='index-wrap'>
        <Text className='text-box head'>
          {protocol.explain}
        </Text>
        <div>
          {
            protocol.rules.map(rule => {
                return (
                  <div>
                    <Text className='level1'>
                      {rule.title}
                    </Text>
                    {rule.content.map(item => {
                      return (
                        <Text className='text-box'>
                          {item}
                        </Text>
                      )
                    })}
                  </div>
                )
              }
            )
          }
        </div>
      </View>
    )
  }
}
