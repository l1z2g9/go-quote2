class AudioTrack extends React.Component {
  constructor(props) {
      super(props);
      this.state = {url: "", title: ""};
      this.handleChange = this.handleChange.bind(this);
      this.playEnd = this.playEnd.bind(this);
   }

  handleChange(e) {
    if (e.target.value.startsWith('http')) {
      this.url = e.target.value;
    }else{
      this.url = "https://drive.google.com/uc?id=" + e.target.value + "&authuser=0&export=download";
    }

    this.title = e.target.dataset.title;

    this.setState({url: this.url, title: this.title});
  }

  playEnd(e){
      let track = Math.floor((Math.random() * this.props.list.length) + 1);
      this.url = "https://drive.google.com/uc?id=" + this.props.list[track][0] + "&authuser=0&export=download";
      this.title = this.props.list[track][1];
      this.setState({url: this.url, title: this.title});
  }

  render() {
    return (
      <div>
         <table className="pure-table">
          <thead>
              <tr className="bg-info">
                  <td>Soundtrack</td>
                  <td>Play</td>
              </tr>
          </thead>
          <tbody>
            {
              this.props.list.map( (list, i) => {
                  return <tr className={i%2==0 ? 'pure-table-odd' : ''} key={list[0]}>
                        <td>
                {list.length > 2 &&
                          
                            <a href="https://drive.google.com/uc?id={list[2]}&authuser=0&export=download">{list[1]}</a>
                }
                {list.length ==2 &&
                            list[1]
                         
                }</td>
                    <td><button type="button" className="pure-button pure-button-primary" data-title={list[1]} onClick={this.handleChange} value={list[0]}>Play</button></td>

                </tr>
             })
            }
          </tbody>
          <tfoot>
              <tr>
                  <td colSpan="5" className="bg-info">
                      <b>{this.state.title}</b>
                  </td>
              </tr>
          </tfoot>
        </table>
        <div>
          <audio controls src={this.state.url} autoPlay onEnded={this.playEnd}></audio>
        </div>
      </div>
    )
  }
}

class Menu extends React.Component {
  constructor(props) {
    super(props);
    this.state = {cat: '粵講越有趣', subCat: '唐詩七絕選賞 - 陳耀南教授主講', list : ycanList};
    this.handleChange = this.handleChange.bind(this);

    this.audioList = {'唐詩七絕選賞 - 陳耀南教授主講':ycanList, '葉培':yePeiList, '潘昭強':panChaoQiangList, '周永傑':zhouYongJieList, 'RTHK': rthkList};
    // this.audioList = {'1':ycanList, '2': yePeiList, '3': panChaoQiangList, '4': zhouYongJieList, '5', zhaobinghengList, '6': liangjieeList, '7': tanhuiqingList, '8': zhongfeiList, '9': yanjianrongList};
  }

  handleChange(e) {
    //e.preventDefault();
    this.cat = e.target.dataset.cat;
    this.subCat = e.target.dataset.subcat;
    this.list = this.audioList[this.subCat];

    this.setState({cat: this.cat, subCat: this.subCat, list: this.list});
  }

  render() {
    return (
      <div>
        <div className="pure-menu pure-menu-horizontal">
            <ul className="pure-menu-list">
                <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="粵講越有趣" data-subcat="唐詩七絕選賞 - 陳耀南教授主講" className="pure-menu-link">粵講越有趣</a></li>
                 <li className="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
                  <a href="#" id="menuLink1" className="pure-menu-link">大城小事</a>
                  <ul className="pure-menu-children">
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="葉培" className="pure-menu-link">葉培</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="潘昭強" className="pure-menu-link">潘昭強</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="周永傑" className="pure-menu-link">周永傑</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="招秉恒" className="pure-menu-link">招秉恒</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="梁潔娥" className="pure-menu-link">梁潔娥</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="譚惠清" className="pure-menu-link">譚惠清</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="鍾飛" className="pure-menu-link">鍾飛</a></li>
                      <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="大城小事" data-subcat="嚴劍蓉" className="pure-menu-link">嚴劍蓉</a></li>
                  </ul>
                </li>
                <li className="pure-menu-item"><a href="#" onClick={this.handleChange} data-cat="LiveRadio" data-subcat="RTHK" className="pure-menu-link">RTHK Live Radio</a></li>
                <li className="pure-menu-item"><a href="#" onClick={this.handleChange} className="pure-menu-link">Finance</a></li>
            </ul>
        </div>
        <h1>{this.state.cat}</h1>
        <h2>{this.state.subCat}</h2>
        <AudioTrack list={this.state.list} />
     </div>
    )
  }
}

ReactDOM.render(<Menu />, document.getElementById('root'));