import Taro, { Component, getStorageSync } from '@tarojs/taro'
import { View, Image, Text } from '@tarojs/components'
import PropTypes from 'prop-types'
import './index.scss'
import img1 from '../../static/icon/img1.png'
import arrow from '../../static/icon/arrow-r.png'
import chapter2 from '../../static/image/chapter2.png'

export default class BorrowItem extends Component {
  constructor() {
    super(...arguments)
		this.state = {
		}
  }
  static defaultProps = {
  		borrowData: {to:{},iou:{},from:{}}
  }
  static propTypes = {
    borrowData: PropTypes.object
	
	}

	componentDidMount(){
		  
	}
	link(id){
		Taro.navigateTo({
			url: '/pages/borrow/borrowPreview/index?iouId='+id
		})		
	}
  render() {
	    let certificateData = getStorageSync('userCertificate');
	    let borrowData = this.props.borrowData;
		  let toOpenId = borrowData.to.OpenId;
		  let OpenId = certificateData.OPENID;
		  let status = borrowData.iou.Status;
		  let borrowName = '借款方';
		  let iou = borrowData.iou;
		  let AvataUrl = '';
		  let NickName = '';
		  let Name = '';
		  let BorrowAt = borrowData.iou.BorrowAt;
		  let PayBackAt = borrowData.iou.PayBackAt;
		  if(status > 1 && OpenId == toOpenId){//当前用户为借钱人并且对方确认了借条
	  			borrowName = '出款方'
	  			AvataUrl = borrowData.from.AvataUrl
	  			NickName = borrowData.from.NickName
	  			Name = borrowData.from.Name
		  }else{
		  		AvataUrl = borrowData.to.AvataUrl
		  		NickName = borrowData.to.NickName
		  		Name = borrowData.to.Name
		  }
			switch(status){
			  case 0:
					status = '待确定'
					break;
				case 1:
					status = '待出借人签名'
					break;
				case 2:
					status = '待出借人签名'
					break;
				case 3:
					status = '待出借人签名'
					break;
				case 4:
					status = '待借款方签名'
					break;
				case 5:
					status = '待出借人放款'
					break;
				case 6:
					status = '出借人已确认放款'
					break;
			  case 7:
					status = '借款方已确认收款'
					break;
			  case 8:
					status = '借款方已确认还款'
					break;	
				case 9:
					status = '借条已结清'
					break;	
				default:
					status = '借条出错'
			}
			let endTime = PayBackAt*1000 //还款日
			let nowTime = new Date().getTime() // 今天
			let nTime = endTime - nowTime
			let day =  Math.floor(nTime / 86400000)+1;
      return (
        <View className='borrow-Item-box' onClick={this.link.bind(this,iou.Id)}>
					<View className='item-head'>
						<View className='user-info'><Text className='txt1'>{borrowName}：</Text><Image className='img' src={AvataUrl} /><Text className='txt2'>{Name}</Text></View>
						<View className='user-state'><Text className='txt'>{status}</Text><Image className='img' src={arrow} /></View>
					</View>
					<View className='item-bottom'>
						<View className='price'>￥{iou.Amount}</View>
						<View className='countdown'>
						{ iou.Status == 9 ? <Image className='chapter-img' src={chapter2} /> : (day <= 0 ? (day < 0 ? <Text className='expire'>借条已逾期{Math.abs(day)}天</Text> :<Text className='expire'>借条今日到期</Text>) : <View><Text className='time'>{day}</Text><Text>天后到期</Text></View>)}
						</View>
					</View>
        </View>
      )
  }
}
