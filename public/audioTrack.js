"use strict";

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var AudioTrack = function (_React$Component) {
  _inherits(AudioTrack, _React$Component);

  function AudioTrack(props) {
    _classCallCheck(this, AudioTrack);

    var _this = _possibleConstructorReturn(this, (AudioTrack.__proto__ || Object.getPrototypeOf(AudioTrack)).call(this, props));

    _this.state = { url: "", title: "" };
    _this.handleChange = _this.handleChange.bind(_this);
    _this.playEnd = _this.playEnd.bind(_this);
    return _this;
  }

  _createClass(AudioTrack, [{
    key: "handleChange",
    value: function handleChange(e) {
      if (e.target.value.startsWith('http')) {
        this.url = e.target.value;
      } else {
        this.url = "https://drive.google.com/uc?id=" + e.target.value + "&authuser=0&export=download";
      }

      this.title = e.target.dataset.title;

      this.setState({ url: this.url, title: this.title });
    }
  }, {
    key: "playEnd",
    value: function playEnd(e) {
      var track = Math.floor(Math.random() * this.props.list.length + 1);
      this.url = "https://drive.google.com/uc?id=" + this.props.list[track][0] + "&authuser=0&export=download";
      this.title = this.props.list[track][1];
      this.setState({ url: this.url, title: this.title });
    }
  }, {
    key: "render",
    value: function render() {
      var _this2 = this;

      return React.createElement(
        "div",
        null,
        React.createElement(
          "table",
          { className: "pure-table" },
          React.createElement(
            "thead",
            null,
            React.createElement(
              "tr",
              { className: "bg-info" },
              React.createElement(
                "td",
                null,
                "Soundtrack"
              ),
              React.createElement(
                "td",
                null,
                "Play"
              )
            )
          ),
          React.createElement(
            "tbody",
            null,
            this.props.list.map(function (list, i) {
              return React.createElement(
                "tr",
                { className: i % 2 == 0 ? 'pure-table-odd' : '', key: list[0] },
                React.createElement(
                  "td",
                  null,
                  list.length > 2 && React.createElement(
                    "a",
                    { href: "https://drive.google.com/uc?id={list[2]}&authuser=0&export=download" },
                    list[1]
                  ),
                  list.length == 2 && list[1]
                ),
                React.createElement(
                  "td",
                  null,
                  React.createElement(
                    "button",
                    { type: "button", className: "pure-button pure-button-primary", "data-title": list[1], onClick: _this2.handleChange, value: list[0] },
                    "Play"
                  )
                )
              );
            })
          ),
          React.createElement(
            "tfoot",
            null,
            React.createElement(
              "tr",
              null,
              React.createElement(
                "td",
                { colSpan: "5", className: "bg-info" },
                React.createElement(
                  "b",
                  null,
                  this.state.title
                )
              )
            )
          )
        ),
        React.createElement(
          "div",
          null,
          React.createElement("audio", { controls: true, src: this.state.url, autoPlay: true, onEnded: this.playEnd })
        )
      );
    }
  }]);

  return AudioTrack;
}(React.Component);

var Menu = function (_React$Component2) {
  _inherits(Menu, _React$Component2);

  function Menu(props) {
    _classCallCheck(this, Menu);

    var _this3 = _possibleConstructorReturn(this, (Menu.__proto__ || Object.getPrototypeOf(Menu)).call(this, props));

    _this3.state = { cat: '粵講越有趣', subCat: '唐詩七絕選賞 - 陳耀南教授主講', list: ycanList };
    _this3.handleChange = _this3.handleChange.bind(_this3);

    _this3.audioList = { '唐詩七絕選賞 - 陳耀南教授主講': ycanList, '葉培': yePeiList, '潘昭強': panChaoQiangList, '周永傑': zhouYongJieList, 'RTHK': rthkList };
    // this.audioList = {'1':ycanList, '2': yePeiList, '3': panChaoQiangList, '4': zhouYongJieList, '5', zhaobinghengList, '6': liangjieeList, '7': tanhuiqingList, '8': zhongfeiList, '9': yanjianrongList};
    return _this3;
  }

  _createClass(Menu, [{
    key: "handleChange",
    value: function handleChange(e) {
      //e.preventDefault();
      this.cat = e.target.dataset.cat;
      this.subCat = e.target.dataset.subcat;
      this.list = this.audioList[this.subCat];

      this.setState({ cat: this.cat, subCat: this.subCat, list: this.list });
    }
  }, {
    key: "render",
    value: function render() {
      return React.createElement(
        "div",
        null,
        React.createElement(
          "div",
          { className: "pure-menu pure-menu-horizontal" },
          React.createElement(
            "ul",
            { className: "pure-menu-list" },
            React.createElement(
              "li",
              { className: "pure-menu-item" },
              React.createElement(
                "a",
                { href: "#", onClick: this.handleChange, "data-cat": "\u7CB5\u8B1B\u8D8A\u6709\u8DA3", "data-subcat": "\u5510\u8A69\u4E03\u7D55\u9078\u8CDE - \u9673\u8000\u5357\u6559\u6388\u4E3B\u8B1B", className: "pure-menu-link" },
                "\u7CB5\u8B1B\u8D8A\u6709\u8DA3"
              )
            ),
            React.createElement(
              "li",
              { className: "pure-menu-item pure-menu-has-children pure-menu-allow-hover" },
              React.createElement(
                "a",
                { href: "#", id: "menuLink1", className: "pure-menu-link" },
                "\u5927\u57CE\u5C0F\u4E8B"
              ),
              React.createElement(
                "ul",
                { className: "pure-menu-children" },
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u8449\u57F9", className: "pure-menu-link" },
                    "\u8449\u57F9"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u6F58\u662D\u5F37", className: "pure-menu-link" },
                    "\u6F58\u662D\u5F37"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u5468\u6C38\u5091", className: "pure-menu-link" },
                    "\u5468\u6C38\u5091"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u62DB\u79C9\u6052", className: "pure-menu-link" },
                    "\u62DB\u79C9\u6052"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u6881\u6F54\u5A25", className: "pure-menu-link" },
                    "\u6881\u6F54\u5A25"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u8B5A\u60E0\u6E05", className: "pure-menu-link" },
                    "\u8B5A\u60E0\u6E05"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u937E\u98DB", className: "pure-menu-link" },
                    "\u937E\u98DB"
                  )
                ),
                React.createElement(
                  "li",
                  { className: "pure-menu-item" },
                  React.createElement(
                    "a",
                    { href: "#", onClick: this.handleChange, "data-cat": "\u5927\u57CE\u5C0F\u4E8B", "data-subcat": "\u56B4\u528D\u84C9", className: "pure-menu-link" },
                    "\u56B4\u528D\u84C9"
                  )
                )
              )
            ),
            React.createElement(
              "li",
              { className: "pure-menu-item" },
              React.createElement(
                "a",
                { href: "#", onClick: this.handleChange, "data-cat": "LiveRadio", "data-subcat": "RTHK", className: "pure-menu-link" },
                "RTHK Live Radio"
              )
            ),
            React.createElement(
              "li",
              { className: "pure-menu-item" },
              React.createElement(
                "a",
                { href: "#", onClick: this.handleChange, className: "pure-menu-link" },
                "Finance"
              )
            )
          )
        ),
        React.createElement(
          "h1",
          null,
          this.state.cat
        ),
        React.createElement(
          "h2",
          null,
          this.state.subCat
        ),
        React.createElement(AudioTrack, { list: this.state.list })
      );
    }
  }]);

  return Menu;
}(React.Component);

ReactDOM.render(React.createElement(Menu, null), document.getElementById('root'));
