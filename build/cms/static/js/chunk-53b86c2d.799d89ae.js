(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-53b86c2d"],{"399c":function(t,e,n){"use strict";n.r(e);var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"app-container"},[n("el-form",{ref:"dataForm",staticStyle:{width:"80%"},attrs:{size:"small","label-width":"100px"}},[t.config?n("el-form-item",{attrs:{label:"录入日期"}},[n("el-date-picker",{staticStyle:{width:"80%"},attrs:{format:"yyyy 年 MM 月 dd 日","value-format":"timestamp",type:"date",placeholder:"Pick a date"},on:{change:t.getData},model:{value:t.day_at,callback:function(e){t.day_at=e},expression:"day_at"}})],1):t._e(),t._v(" "),t.config?n("el-form-item",{attrs:{label:"微剧列表",prop:"name"}},[t._l(t.detail,(function(e,a){return n("div",{key:a},[n("video-item",{attrs:{item:e,config:t.config.platform_micro},on:{del:function(e){return t.detail.splice(a,1)},"update:item":function(t){e=t}}})],1)})),t._v(" "),n("el-button",{staticStyle:{width:"80%"},on:{click:function(e){return t.detail.push({name:"",platform:[],num:null})}}},[t._v("\n        +新增\n      ")])],2):t._e(),t._v(" "),n("el-form-item",[n("el-button",{attrs:{type:"primary"},on:{click:t.createData}},[t._v("\n        保存\n      ")])],1)],1)],1)},i=[],l=n("b775"),r=n("7590"),s=n("bcec"),c=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticStyle:{"margin-bottom":"5px"}},[n("el-row",[n("el-col",{staticStyle:{"padding-right":"10px"},attrs:{span:8}},[n("el-input",{attrs:{placeholder:"名称"},model:{value:t.item.name,callback:function(e){t.$set(t.item,"name",e)},expression:"item.name"}})],1),t._v(" "),n("el-col",{staticStyle:{"padding-right":"10px"},attrs:{span:6}},[n("el-select",{staticStyle:{width:"100%"},attrs:{multiple:"",placeholder:"请选择"},model:{value:t.item.platform,callback:function(e){t.$set(t.item,"platform",e)},expression:"item.platform"}},t._l(t.config,(function(t,e){return n("el-option",{key:e,attrs:{label:t.name,value:t.id}})})),1)],1),t._v(" "),n("el-col",{staticStyle:{"padding-right":"10px"},attrs:{span:4}},[n("el-input-number",{staticStyle:{width:"100%"},attrs:{controls:!1,placeholder:"分值",min:0,max:100},model:{value:t.item.num,callback:function(e){t.$set(t.item,"num",e)},expression:"item.num"}})],1),t._v(" "),n("el-col",{staticStyle:{"text-align":"left"},attrs:{span:4}},[n("el-button",{attrs:{icon:"el-icon-minus"},on:{click:function(e){return t.$emit("del")}}})],1)],1)],1)},o=[],u=(n("96cf"),n("3b8d")),m={name:"staffItem",props:{item:{type:Object,default:function(){}},config:{type:Array,default:function(){return[]}}},data:function(){return{}},methods:{del:function(){var t=this;this.$confirm("确认删除?","Warning",{confirmButtonText:"删除",cancelButtonText:"取消",type:"warning"}).then(Object(u["a"])(regeneratorRuntime.mark((function e(){return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:t.$emit("del");case 1:case"end":return e.stop()}}),e)})))).catch((function(t){console.error(t)}))}}},p=m,f=n("2877"),d=Object(f["a"])(p,c,o,!1,null,"5a30aa96",null),h=d.exports,y=n("4ec3"),b={name:"name",components:{staffItem:r["a"],actorItem:s["a"],videoItem:h},data:function(){return{host:l["a"],show_type:4,detail:{},day_at:new Date((new Date).toDateString()).getTime(),staff_type:[{name:"编剧",id:"screen_writer"},{name:"制片人",id:"producer"},{name:"制片公司",id:"producer_company"},{name:"出品人",id:"publisher"},{name:"出品公司",id:"publisher_company"}],rules:{role_id:[{required:!0,message:"type is required",trigger:"change"}],timestamp:[{type:"date",required:!0,message:"timestamp is required",trigger:"change"}],email:[{required:!0,message:"title is required",trigger:"blur"}]},config:null}},mounted:function(){this.getConfig(),this.getData()},methods:{getConfig:function(){var t=this;Object(y["b"])().then((function(e){t.config=e.data}))},pushStaff:function(t,e){},getData:function(){var t=this;Object(y["c"])({day_at:this.day_at/1e3}).then((function(e){t.detail=e.data}))},createData:function(){var t=this;console.log(this.detail);var e={};e.data=this.detail,e.day_at=this.day_at?this.day_at/1e3:0,e=JSON.stringify(e),Object(y["d"])(e).then((function(e){t.$notify({title:"Success",type:"success",duration:2e3})}))}}},g=b,v=(n("e6bd"),Object(f["a"])(g,a,i,!1,null,null,null));e["default"]=v.exports},4603:function(t,e,n){},7590:function(t,e,n){"use strict";var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",[t._l(t.staffs,(function(e){return n("el-tag",{key:e,attrs:{type:t.theme,closable:"","disable-transitions":!1},on:{close:function(n){return t.handleClose(e)}}},[t._v("\n    "+t._s(e)+"\n  ")])})),t._v(" "),t.inputVisible?n("el-input",{ref:"saveTagInput",staticClass:"input-new-tag",attrs:{size:"small"},on:{blur:t.handleInputConfirm},nativeOn:{keyup:function(e){return!e.type.indexOf("key")&&t._k(e.keyCode,"enter",13,e.key,"Enter")?null:t.handleInputConfirm(e)}},model:{value:t.inputValue,callback:function(e){t.inputValue=e},expression:"inputValue"}}):n("el-button",{staticClass:"button-new-tag",attrs:{size:"small"},on:{click:t.showInput}},[t._v("+ 新增")])],2)},i=[],l={props:{theme:{type:String,default:"primary"},staffs:{type:Array,default:function(){return[]}}},data:function(){return{inputVisible:!1,inputValue:""}},methods:{handleClose:function(t){this.$emit("del",this.staffs.indexOf(t))},showInput:function(){var t=this;this.inputVisible=!0,this.$nextTick((function(e){t.$refs.saveTagInput.$refs.input.focus()}))},handleInputConfirm:function(){var t=this.inputValue;t&&this.$emit("add",t),this.inputVisible=!1,this.inputValue=""}}},r=l,s=(n("85f8"),n("2877")),c=Object(s["a"])(r,a,i,!1,null,"571d98e4",null);e["a"]=c.exports},"85b1":function(t,e,n){},"85f8":function(t,e,n){"use strict";var a=n("4603"),i=n.n(a);i.a},bcec:function(t,e,n){"use strict";var a=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticStyle:{"margin-bottom":"5px"}},[n("span",[t._v("姓名:")]),t._v(" "),n("el-input",{staticStyle:{width:"150px","margin-right":"10px"},model:{value:t.item.name,callback:function(e){t.$set(t.item,"name",e)},expression:"item.name"}}),t._v(" "),n("span",[t._v("饰演:")]),t._v(" "),n("el-input",{staticStyle:{width:"150px","margin-right":"10px"},model:{value:t.item.play,callback:function(e){t.$set(t.item,"play",e)},expression:"item.play"}}),t._v(" "),n("el-select",{attrs:{placeholder:"请选择"},model:{value:t.item.play_type,callback:function(e){t.$set(t.item,"play_type",e)},expression:"item.play_type"}},t._l(t.config,(function(t,e){return n("el-option",{key:e,attrs:{label:t.name,value:t.id}})})),1),t._v(" "),n("el-button-group",[n("el-button",{attrs:{icon:"el-icon-plus"},on:{click:function(e){return t.$emit("add")}}}),t._v(" "),n("el-button",{attrs:{icon:"el-icon-minus"},on:{click:function(e){return t.$emit("del")}}})],1)],1)},i=[],l={name:"staffItem",props:{item:{type:Object,default:function(){}},config:{type:Array,default:function(){return[]}}},data:function(){return{}},methods:{}},r=l,s=n("2877"),c=Object(s["a"])(r,a,i,!1,null,"6bea7400",null);e["a"]=c.exports},e6bd:function(t,e,n){"use strict";var a=n("85b1"),i=n.n(a);i.a}}]);