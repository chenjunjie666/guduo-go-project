(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-39f8586a"],{"2b91":function(t,e,a){},"8c01":function(t,e,a){"use strict";var n=a("2b91"),s=a.n(n);s.a},fa7a:function(t,e,a){"use strict";a.r(e);var n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"app-container"},[a("div",{staticStyle:{"margin-bottom":"10px",display:"flex","align-items":"center"}},[a("span",{staticStyle:{"font-size":"14px",color:"#333","line-height":"14px"}},[t._v("关键字：")]),t._v(" "),a("el-input",{staticStyle:{width:"30%"},attrs:{size:"small",placeholder:"请输入影视名称"},model:{value:t.searchParams.keyword,callback:function(e){t.$set(t.searchParams,"keyword",e)},expression:"searchParams.keyword"}}),t._v(" "),a("el-button",{staticStyle:{margin:"0 10px"},attrs:{size:"small",type:"primary"},on:{click:t.getList}},[t._v("搜索")])],1),t._v(" "),a("div",{staticStyle:{"margin-bottom":"10px"}},[a("el-button",{attrs:{size:"mini",type:"success",icon:"el-icon-plus"},on:{click:t.handleCreate}},[t._v("添加")])],1),t._v(" "),a("el-table",{directives:[{name:"loading",rawName:"v-loading",value:t.listLoading,expression:"listLoading"}],attrs:{data:t.list,"ma-height":t.tableHeight,"row-style":{height:"40px"},"cell-style":{padding:"0"},"element-loading-text":"Loading",border:"",fit:"","highlight-current-row":""}},[a("el-table-column",{attrs:{type:"index",width:"50"}}),t._v(" "),a("el-table-column",{attrs:{label:"综艺名",width:"150"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{staticStyle:{"white-space":"initial"}},[t._v("\n        "+t._s(e.row.name)+"\n      ")])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"类型"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{},[t._v("\n        "+t._s(e.row.type_str)+"\n      ")])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"上映时间"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{},[e.row.release_at?a("span",[t._v(t._s(t._f("parseTime")(e.row.release_at,"{y}-{m}-{d}")))]):t._e()])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"放映状态"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{},[t._v("\n        "+t._s(e.row.show_status_str)+"\n      ")])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"是否显示"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{},[t._v("\n        "+t._s(e.row.is_show_str)+"\n      ")])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"状态"},scopedSlots:t._u([{key:"default",fn:function(e){return a("div",{},[t._v("\n        "+t._s(e.row.status_str)+"\n      ")])}}])}),t._v(" "),a("el-table-column",{attrs:{label:"操作",width:"160","class-name":"small-padding fixed-width"},scopedSlots:t._u([{key:"default",fn:function(e){var n=e.row,s=e.$index;return a("div",{},[a("router-link",{attrs:{to:"/video-list/variety-show-edit/"+n.show_id}},[a("el-button",{attrs:{type:"primary",size:"mini"}},[t._v("\n            编辑\n          ")])],1),t._v(" "),a("el-button",{attrs:{type:"danger",size:"mini"},on:{click:function(e){return t.handleDelete(n,s)}}},[t._v("\n          删除\n        ")])],1)}}])})],1),t._v(" "),a("el-pagination",{staticStyle:{float:"right","margin-top":"10px"},attrs:{background:"",layout:"prev, pager, next","current-page":t.listParams.page,"page-size":t.listParams.limit,total:t.total},on:{"current-change":t.pageChange}})],1)},s=[],i=(a("96cf"),a("3b8d")),l=a("db72"),r=a("4ec3"),o=a("ed08"),c={filters:{statusFilter:function(t){var e={published:"success",draft:"gray",deleted:"danger"};return e[t]},parseTime:o["b"]},data:function(){return{searchParams:{keyword:"",type:""},listParams:{page:1,limit:10},list:null,listLoading:!0,tableHeight:window.innerHeight-200,total:0}},created:function(){this.getList()},methods:{pageChange:function(t){this.listParams.page=t,this.getList()},handleCreate:function(){this.$router.push({name:"variety-show-create"})},getList:function(){var t=this;this.listLoading=!0,Object(r["j"])(Object(l["a"])({},this.searchParams,{},this.listParams,{show_type:1})).then((function(e){console.log(e),t.list=e.data.data,t.total=e.data.total,t.listLoading=!1}))},handleDelete:function(t,e){var a=this;this.$confirm("确认删除?","警告",{confirmButtonText:"删除",cancelButtonText:"取消",type:"warning"}).then(Object(i["a"])(regeneratorRuntime.mark((function n(){return regeneratorRuntime.wrap((function(n){while(1)switch(n.prev=n.next){case 0:return n.next=2,Object(r["g"])({show_id:t.show_id});case 2:a.list.splice(e,1),a.$message({type:"success",message:"Delete Successfully!"});case 4:case"end":return n.stop()}}),n)})))).catch((function(t){console.error(t)}))}}},u=c,d=(a("8c01"),a("2877")),p=Object(d["a"])(u,n,s,!1,null,null,null);e["default"]=p.exports}}]);