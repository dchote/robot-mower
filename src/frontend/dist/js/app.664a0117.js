(function(t){function e(e){for(var a,l,s=e[0],i=e[1],c=e[2],f=0,v=[];f<s.length;f++)l=s[f],r[l]&&v.push(r[l][0]),r[l]=0;for(a in i)Object.prototype.hasOwnProperty.call(i,a)&&(t[a]=i[a]);u&&u(e);while(v.length)v.shift()();return o.push.apply(o,c||[]),n()}function n(){for(var t,e=0;e<o.length;e++){for(var n=o[e],a=!0,s=1;s<n.length;s++){var i=n[s];0!==r[i]&&(a=!1)}a&&(o.splice(e--,1),t=l(l.s=n[0]))}return t}var a={},r={app:0},o=[];function l(e){if(a[e])return a[e].exports;var n=a[e]={i:e,l:!1,exports:{}};return t[e].call(n.exports,n,n.exports,l),n.l=!0,n.exports}l.m=t,l.c=a,l.d=function(t,e,n){l.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:n})},l.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},l.t=function(t,e){if(1&e&&(t=l(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var n=Object.create(null);if(l.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var a in t)l.d(n,a,function(e){return t[e]}.bind(null,a));return n},l.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return l.d(e,"a",e),e},l.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},l.p="/";var s=window["webpackJsonp"]=window["webpackJsonp"]||[],i=s.push.bind(s);s.push=e,s=s.slice();for(var c=0;c<s.length;c++)e(s[c]);var u=i;o.push([0,"chunk-vendors"]),n()})({0:function(t,e,n){t.exports=n("56d7")},"034f":function(t,e,n){"use strict";var a=n("c21b"),r=n.n(a);r.a},"2a83":function(t,e,n){},"2e4a":function(t,e,n){"use strict";var a=n("3516"),r=n.n(a);r.a},3516:function(t,e,n){},"35da":function(t,e,n){"use strict";var a=n("e5b9"),r=n.n(a);r.a},"56d7":function(t,e,n){"use strict";n.r(e);n("744f"),n("6c7b"),n("7514"),n("20d6"),n("1c4c"),n("6762"),n("cadf"),n("e804"),n("55dd"),n("d04f"),n("c8ce"),n("217b"),n("7f7f"),n("f400"),n("7f25"),n("536b"),n("d9ab"),n("f9ab"),n("32d7"),n("25c9"),n("9f3c"),n("042e"),n("c7c6"),n("f4ff"),n("049f"),n("7872"),n("a69f"),n("0b21"),n("6c1a"),n("c7c62"),n("84b4"),n("c5f6"),n("2e37"),n("fca0"),n("7cdf"),n("ee1d"),n("b1b1"),n("87f3"),n("9278"),n("5df2"),n("04ff"),n("f751"),n("4504"),n("fee7"),n("ffc1"),n("0d6d"),n("9986"),n("8e6e"),n("25db"),n("e4f7"),n("b9a1"),n("64d5"),n("9aea"),n("db97"),n("66c8"),n("57f0"),n("165b"),n("456d"),n("cf6a"),n("fd24"),n("8615"),n("551c"),n("097d"),n("df1b"),n("2397"),n("88ca"),n("ba16"),n("d185"),n("ebde"),n("2d34"),n("f6b3"),n("2251"),n("c698"),n("a19f"),n("9253"),n("9275"),n("3b2b"),n("3846"),n("4917"),n("a481"),n("28a5"),n("386d"),n("6b54"),n("4f7f"),n("8a81"),n("ac4d"),n("8449"),n("9c86"),n("fa83"),n("48c0"),n("a032"),n("aef6"),n("d263"),n("6c37"),n("9ec8"),n("5695"),n("2fdb"),n("d0b0"),n("b54a"),n("f576"),n("ed50"),n("788d"),n("14b9"),n("f386"),n("f559"),n("1448"),n("673e"),n("242a"),n("c66f"),n("b05c"),n("34ef"),n("6aa2"),n("15ac"),n("af56"),n("b6e4"),n("9c29"),n("63d9"),n("4dda"),n("10ad"),n("c02b"),n("4795"),n("130f"),n("ac6a"),n("96cf");var a=n("2b0e"),r=n("ce5b"),o=n.n(r);n("bf40");a["default"].use(o.a,{});var l=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("v-app",{attrs:{id:"mower"}},[n("v-navigation-drawer",{attrs:{clipped:t.$vuetify.breakpoint.lgAndUp,fixed:"",app:""},model:{value:t.drawer,callback:function(e){t.drawer=e},expression:"drawer"}},[n("v-list",{attrs:{dense:""}},[t._l(t.items,function(e){return[e.heading?n("v-layout",{key:e.heading,attrs:{row:"","align-center":""}},[n("v-flex",{attrs:{xs6:""}},[e.heading?n("v-subheader",[t._v("\n              "+t._s(e.heading)+"\n            ")]):t._e()],1),n("v-flex",{staticClass:"text-xs-center",attrs:{xs6:""}},[n("a",{staticClass:"body-2 black--text",attrs:{href:"#!"}},[t._v("EDIT")])])],1):e.children?n("v-list-group",{key:e.text,attrs:{"prepend-icon":e.model?e.icon:e["icon-alt"],"append-icon":""},model:{value:e.model,callback:function(n){t.$set(e,"model",n)},expression:"item.model"}},[n("v-list-tile",{attrs:{slot:"activator"},slot:"activator"},[n("v-list-tile-content",[n("v-list-tile-title",[t._v("\n                "+t._s(e.text)+"\n              ")])],1)],1),t._l(e.children,function(e,a){return n("v-list-tile",{key:a,on:{click:function(n){t.$router.push(e.route)}}},[e.icon?n("v-list-tile-action",[n("v-icon",[t._v(t._s(e.icon))])],1):t._e(),n("v-list-tile-content",[n("v-list-tile-title",[t._v("\n                "+t._s(e.text)+"\n              ")])],1)],1)})],2):n("v-list-tile",{key:e.text,on:{click:function(n){t.$router.push(e.route)}}},[n("v-list-tile-action",[n("v-icon",[t._v(t._s(e.icon))])],1),n("v-list-tile-content",[n("v-list-tile-title",[t._v("\n              "+t._s(e.text)+"\n            ")])],1)],1)]})],2)],1),n("v-toolbar",{attrs:{"clipped-left":t.$vuetify.breakpoint.lgAndUp,color:"blue darken-4",dark:"",app:"",fixed:""}},[n("v-toolbar-title",{staticClass:"ml-0 pl-3",staticStyle:{width:"300px"}},[n("v-toolbar-side-icon",{on:{click:function(e){e.stopPropagation(),t.drawer=!t.drawer}}}),n("span",{staticClass:"hidden-sm-and-down"},[t._v("Robot Mower")])],1),n("v-spacer"),n("v-btn",{attrs:{icon:""}},[n("v-icon",[t._v("notifications")])],1)],1),n("v-content",{staticClass:"cameraView"},[n("v-container",{attrs:{"fill-height":""}},[n("router-view")],1),n("status-bar")],1)],1)},s=[],i={data:function(){return{drawer:null,items:[{icon:"control_camera",text:"Control",route:"control"},{icon:"map",text:"Planner",route:"planner"},{icon:"settings",text:"Settings",route:"settings"}]}},props:{source:String}},c=i,u=(n("034f"),n("2877")),f=Object(u["a"])(c,l,s,!1,null,null,null);f.options.__file="App.vue";var v=f.exports,d=n("8c4f"),p=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("v-layout",{attrs:{"align-center":"","justify-center":"",row:"","fill-height":""}},[a("section",[a("v-layout",{attrs:{column:"",wrap:"","align-center":"",grey:"","lighten-4":"","elevation-4":""}},[a("v-flex",{staticClass:"my-3",attrs:{xs12:"",sm4:"","pa-5":""}},[a("div",{staticClass:"text-xs-center"},[a("h2",{staticClass:"headline"},[t._v("Welcome to your Robot Mower")]),a("span",{staticClass:"subheading"},[t._v("\n              To control your mower, please make a selection from the menu on the left.\n            ")]),a("img",{staticClass:"mt-4",attrs:{src:n("7c01"),width:"70%"}})])])],1)],1)])},b=[],m={name:"Welcome",data:function(){return{}}},_=m,h=(n("2e4a"),Object(u["a"])(_,p,b,!1,null,"6f96cc75",null));h.options.__file="Welcome.vue";var g=h.exports,x=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("v-layout",{attrs:{"align-end":"","justify-center":"",row:"","fill-height":""}},[n("v-card",{attrs:{height:"70px",width:"700px",flat:""}},[n("v-bottom-nav",{attrs:{absolute:"",color:"transparent",value:!0}},[n("v-flex",{attrs:{xs3:"","mr-1":"","ml-4":"","mt-2":""}},[n("v-slider",{attrs:{"inverse-label":"",label:"Speed"},model:{value:t.speed,callback:function(e){t.speed=e},expression:"speed"}})],1),n("v-layout",{attrs:{xs3:""}},[n("v-flex",{attrs:{xs3:""}},[n("v-btn",{attrs:{color:"teal",flat:""}},[n("span",[t._v("Left")]),n("v-icon",[t._v("arrow_back")])],1)],1),n("v-flex",{attrs:{xs3:""}},[n("v-btn",{attrs:{color:"teal",flat:""}},[n("span",[t._v("Forward")]),n("v-icon",[t._v("arrow_upward")])],1)],1),n("v-flex",{attrs:{xs3:""}},[n("v-btn",{attrs:{color:"teal",flat:""}},[n("span",[t._v("Backward")]),n("v-icon",[t._v("arrow_downward")])],1)],1),n("v-flex",{attrs:{xs3:""}},[n("v-btn",{attrs:{color:"teal",flat:""}},[n("span",[t._v("Right")]),n("v-icon",[t._v("arrow_forward")])],1)],1)],1),n("v-flex",{attrs:{xs3:"","ml-1":"","mr-4":"","mt-2":""}},[n("v-slider",{attrs:{label:"Cutter"},model:{value:t.cutter,callback:function(e){t.cutter=e},expression:"cutter"}})],1)],1)],1)],1)},w=[],y={name:"ControlPage",data:function(){return{speed:100,cutter:75}}},k=y,C=(n("35da"),Object(u["a"])(k,x,w,!1,null,"29a3b8f9",null));C.options.__file="Control.vue";var j=C.exports,O=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("v-container",{attrs:{fluid:""}},[n("v-slide-y-transition",{attrs:{mode:"out-in"}},[n("v-layout",{attrs:{column:"","align-center":""}},[n("v-alert",{attrs:{value:!0,color:"error",icon:"warning",outline:""}},[t._v("\n           Not yet implemented.\n         ")])],1)],1)],1)},S=[],P={},$=P,E=(n("ba89"),Object(u["a"])($,O,S,!1,null,"18a44104",null));E.options.__file="NotImplemented.vue";var M=E.exports,T=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"statsContainer"},[n("v-layout",{attrs:{"align-center":"","justify-end":"",row:"","fill-height":""}},[n("div",{staticClass:"stat grey lighten-4 elevation-2 text-xs-center"},[n("h5",[t._v("Battery Voltage:")]),t._v("\n      "+t._s(t.voltage)+"\n    ")]),n("div",{staticClass:"stat grey lighten-4 elevation-2 text-xs-center"},[n("h5",[t._v("Compass:")]),t._v("\n      "+t._s(t.compass)+"\n    ")]),n("div",{staticClass:"stat grey lighten-4 elevation-2 text-xs-center"},[n("h5",[t._v("GPS:")]),t._v("\n      "+t._s(t.gps)+"\n    ")])])],1)},B=[],W={name:"StatusBar",data:function(){return{voltage:"24.1",compass:"NE",gps:"40.780715, -78.007729"}}},A=W,N=(n("d4aa"),Object(u["a"])(A,T,B,!1,null,"5c850341",null));N.options.__file="StatusBar.vue";var R=N.exports;a["default"].use(d["a"]),a["default"].component("status-bar",R);var I=new d["a"]({routes:[{path:"/",name:"Welcome",component:g},{path:"/control",name:"Control",component:j},{path:"/planner",name:"Planner",component:M},{path:"/settings",name:"Settings",component:M}]}),J=n("2f62");a["default"].use(J["a"]);var U={},V={},D={},F={},G=new J["a"].Store({state:U,getters:F,actions:D,mutations:V});a["default"].config.productionTip=!1,new a["default"]({router:I,store:G,render:function(t){return t(v)}}).$mount("#app")},"75de":function(t,e,n){},"7c01":function(t,e,n){t.exports=n.p+"img/Mower_Chassis_v35.2bd12d27.png"},ba89:function(t,e,n){"use strict";var a=n("2a83"),r=n.n(a);r.a},c21b:function(t,e,n){},d4aa:function(t,e,n){"use strict";var a=n("75de"),r=n.n(a);r.a},e5b9:function(t,e,n){}});
//# sourceMappingURL=app.664a0117.js.map