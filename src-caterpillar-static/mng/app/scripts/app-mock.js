/*global $:false */
'use strict';

(function () {
  $('body[ng-app]').attr('ng-app', function() {
    return $(this).attr('ng-app') + 'Dev';
  });

  // モックappを作成して登録（名前は対象app名 + 'Dev'）
  var appDev = angular.module('caterpillarAppDev', ['caterpillarApp', 'ngMockE2E'/*, 'mockCommon'*/]);
  appDev.run(['$httpBackend', function($httpBackend) {

  //////////////////////////// Page ///////////////////////////////////
  
  // 正規表現にマッチするパスに対してXHRがあったらモックレスポンス(JSON)を返却
  $httpBackend.whenPUT(/\/caterpillar\/api\/pages\/?.*/).respond({result: 'success'});
  $httpBackend.whenPOST(/\/caterpillar\/api\/pages\/?.*/).respond({result: 'success'});
 
  $httpBackend.whenGET('/caterpillar/api/pages').respond([
    {
      key: 'key-1',
      id : 1,
      name : 'admin template',
      alias : 'tmpl_admin',
      leaf : 'leaf_01',
      viewURL : '/1.view',
      editURL : '/1.view',
      createdAt : '2014-07-25T10:03:04.587579Z',
      updatedAt : '2014-07-25T10:03:04.587579Z',
      properties : { 'wh_01_01':'orange', 'wh_01_02':'apple', 'wh_01_03':'banana', 'wh_01_05':'cherry' }
    },
    {
      key: 'key-2',
      id : 2,
      name : 'web template',
      alias : 'tmpl_web',
      leaf : 'leaf_02',
      viewURL : '/2.view',
      editURL : '/2.view',
      createdAt : '2014-07-25T10:03:04.587579Z',
      updatedAt : '2014-07-25T10:03:04.587579Z',
      properties : { 'wh_02_02':'mary', 'wh_02_03':'tony', 'wh_02_05':'jerry' }
    },
    {
      key: 'key-3',
      id : 3,
      name : 'phone template',
      alias : 'tmpl_phone',
      leaf : 'leaf_03',
      viewURL : '/3.view',
      editURL : '/3.view',
      createdAt : '2014-07-25T10:03:04.587579Z',
      updatedAt : '2014-07-25T10:03:04.587579Z',
      properties : { 'wh_03_01':'golang', 'wh_03_02':'martini', 'wh_03_03':'gae', 'wh_03_04':'angular', 'wh_03_05':'jquery' }
    }
  ]);

  $httpBackend.whenGET('/caterpillar/api/leaves').respond(
    [
      {
        name : 'leaf_01',
        alias : 'lf1',
        wormholes : [
          { name:'Global_a', alias:'ga', global:true, type:'PROPERTY' },
          { name:'wh_01_02', alias:'wh02', global:false, type:'PROPERTY' },
          { name:'wh_01_03', global:false, type:'PROPERTY' },
          { name:'Global_b', alias:'', global:true, type:'AREA' },
          { name:'Global_c', alias:'gc', global:true, type:'PROPERTY' }
        ]
      },
      {
        name : 'leaf_02',
        alias : 'lf2',
        wormholes : [
          { name:'Global_b', alias:'', global:true, type:'AREA' },
          { name:'wh_02_02', alias:'', global:false, type:'PROPERTY' },
          { name:'wh_02_03', alias:'wh08', global:false, type:'PROPERTY' },
          { name:'Global_c', alias:'gc', global:true, type:'PROPERTY' },
          { name:'wh_02_05', alias:'wh10', global:false, type:'PROPERTY' }
        ]
      },
      {
        name : 'leaf_03',
        alias : 'lf3',
        wormholes : [
          { name:'Global_a', alias:'ga', global:true, type:'PROPERTY' },
          { name:'wh_03_02', alias:'wh12', global:false, type:'PROPERTY' },
          { name:'Global_b', alias:'', global:true, type:'PROPERTY' },
          { name:'wh_03_04', alias:'', global:false, type:'PROPERTY' },
          { name:'wh_03_05', alias:'wh15', global:false, type:'PROPERTY' }
        ]
      }
    ]
  );

  $httpBackend.whenGET(/\/caterpillar\/api\/pages\/.+/).respond(function(method, url) {
    var id = url.substring(url.lastIndexOf('/')+1);

    var response;
      if (id === '1') {
        response = {
          key: 'key-1',
          id : 1,
          name : 'admin template',
          alias : 'tmpl_admin',
          leaf : 'leaf_01',
          viewURL : '/1.view',
          editURL : '/1.view',
          createdAt : '2014-07-25T10:03:04.587579Z',
          updatedAt : '2014-07-25T10:03:04.587579Z',
          properties : { 'wh_01_01':'orange', 'wh_01_02':'apple', 'wh_01_03':'banana', 'wh_01_05':'cherry' }
        };
      } else if (id === '2') {
        response = {
          key: 'key-2',
          id : 2,
          name : 'web template',
          alias : 'tmpl_web',
          leaf : 'leaf_02',
          viewURL : '/2.view',
          editURL : '/2.view',
          createdAt : '2014-07-25T10:03:04.587579Z',
          updatedAt : '2014-07-25T10:03:04.587579Z',
          properties : { 'wh_02_02':'mary', 'wh_02_03':'tony', 'wh_02_05':'jerry' }
        };
      } else {
        response = {
          key: 'key-3',
          id : 3,
          name : 'phone template',
          alias : 'tmpl_phone',
          leaf : 'leaf_03',
          viewURL : '/3.view',
          editURL : '/3.view',
          createdAt : '2014-07-25T10:03:04.587579Z',
          updatedAt : '2014-07-25T10:03:04.587579Z',
          properties : { 'wh_03_01':'golang', 'wh_03_02':'martini', 'wh_03_03':'gae', 'wh_03_04':'angular', 'wh_03_05':'jquery' }
        };
      }

    return [200, response];
  });

  $httpBackend.whenGET(/\/caterpillar\/api\/property\/.+/).respond(function(method, url) {
    var name = url.substring(url.lastIndexOf('/')+1);

    var response;
    if (name === 'Global_a') {
      response = {value:'property_value_apple'};
    } else if (name === 'Global_b') {
      response = {value:'property_value_banana'};
    } else if (name === 'Global_c') {
      response = {value:'property_value_car'};
    } else {
      response = {value:'property_value'};
    }

    return [200, response];
  });
  
  $httpBackend.whenPUT('/caterpillar/api/rootPage').respond({});

  // htmlファイルの取得等はそのままスルー
  $httpBackend.whenGET(/views\/.*/).passThrough();

  }]);
  
}());
