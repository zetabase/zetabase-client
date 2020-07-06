/**
 * @fileoverview gRPC-Web generated client stub for zbprotocol
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.zbprotocol = require('./zbprotocol_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.zbprotocol.ZetabaseProviderClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.zbprotocol.ZetabaseProviderPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.ZbEmpty,
 *   !proto.zbprotocol.VersionDetails>}
 */
const methodDescriptor_ZetabaseProvider_VersionInfo = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/VersionInfo',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.ZbEmpty,
  proto.zbprotocol.VersionDetails,
  /**
   * @param {!proto.zbprotocol.ZbEmpty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.VersionDetails.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.ZbEmpty,
 *   !proto.zbprotocol.VersionDetails>}
 */
const methodInfo_ZetabaseProvider_VersionInfo = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.VersionDetails,
  /**
   * @param {!proto.zbprotocol.ZbEmpty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.VersionDetails.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.ZbEmpty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.VersionDetails)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.VersionDetails>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.versionInfo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/VersionInfo',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_VersionInfo,
      callback);
};


/**
 * @param {!proto.zbprotocol.ZbEmpty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.VersionDetails>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.versionInfo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/VersionInfo',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_VersionInfo);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.SubIdentityModify,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_ModifySubIdentity = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/ModifySubIdentity',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.SubIdentityModify,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.SubIdentityModify} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.SubIdentityModify,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_ModifySubIdentity = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.SubIdentityModify} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.SubIdentityModify} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.modifySubIdentity =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ModifySubIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ModifySubIdentity,
      callback);
};


/**
 * @param {!proto.zbprotocol.SubIdentityModify} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.modifySubIdentity =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ModifySubIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ModifySubIdentity);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.SimpleRequest,
 *   !proto.zbprotocol.SubIdentitiesList>}
 */
const methodDescriptor_ZetabaseProvider_ListSubIdentities = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/ListSubIdentities',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.SimpleRequest,
  proto.zbprotocol.SubIdentitiesList,
  /**
   * @param {!proto.zbprotocol.SimpleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.SubIdentitiesList.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.SimpleRequest,
 *   !proto.zbprotocol.SubIdentitiesList>}
 */
const methodInfo_ZetabaseProvider_ListSubIdentities = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.SubIdentitiesList,
  /**
   * @param {!proto.zbprotocol.SimpleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.SubIdentitiesList.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.SimpleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.SubIdentitiesList)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.SubIdentitiesList>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.listSubIdentities =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListSubIdentities',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListSubIdentities,
      callback);
};


/**
 * @param {!proto.zbprotocol.SimpleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.SubIdentitiesList>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.listSubIdentities =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListSubIdentities',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListSubIdentities);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.NewIdentityRequest,
 *   !proto.zbprotocol.NewIdentityResponse>}
 */
const methodDescriptor_ZetabaseProvider_RegisterNewIdentity = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/RegisterNewIdentity',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.NewIdentityRequest,
  proto.zbprotocol.NewIdentityResponse,
  /**
   * @param {!proto.zbprotocol.NewIdentityRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.NewIdentityResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.NewIdentityRequest,
 *   !proto.zbprotocol.NewIdentityResponse>}
 */
const methodInfo_ZetabaseProvider_RegisterNewIdentity = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.NewIdentityResponse,
  /**
   * @param {!proto.zbprotocol.NewIdentityRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.NewIdentityResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.NewIdentityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.NewIdentityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.NewIdentityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.registerNewIdentity =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/RegisterNewIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_RegisterNewIdentity,
      callback);
};


/**
 * @param {!proto.zbprotocol.NewIdentityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.NewIdentityResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.registerNewIdentity =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/RegisterNewIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_RegisterNewIdentity);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.NewIdentityConfirm,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_ConfirmNewIdentity = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/ConfirmNewIdentity',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.NewIdentityConfirm,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.NewIdentityConfirm} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.NewIdentityConfirm,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_ConfirmNewIdentity = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.NewIdentityConfirm} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.NewIdentityConfirm} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.confirmNewIdentity =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ConfirmNewIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ConfirmNewIdentity,
      callback);
};


/**
 * @param {!proto.zbprotocol.NewIdentityConfirm} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.confirmNewIdentity =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ConfirmNewIdentity',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ConfirmNewIdentity);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.TablePut,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_PutData = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/PutData',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.TablePut,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TablePut} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.TablePut,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_PutData = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TablePut} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.TablePut} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.putData =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/PutData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_PutData,
      callback);
};


/**
 * @param {!proto.zbprotocol.TablePut} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.putData =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/PutData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_PutData);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.TablePutMulti,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_PutDataMulti = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/PutDataMulti',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.TablePutMulti,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TablePutMulti} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.TablePutMulti,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_PutDataMulti = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TablePutMulti} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.TablePutMulti} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.putDataMulti =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/PutDataMulti',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_PutDataMulti,
      callback);
};


/**
 * @param {!proto.zbprotocol.TablePutMulti} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.putDataMulti =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/PutDataMulti',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_PutDataMulti);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.TableCreate,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_CreateTable = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/CreateTable',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.TableCreate,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TableCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.TableCreate,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_CreateTable = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.TableCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.TableCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.createTable =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/CreateTable',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_CreateTable,
      callback);
};


/**
 * @param {!proto.zbprotocol.TableCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.createTable =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/CreateTable',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_CreateTable);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.TableGet,
 *   !proto.zbprotocol.TableGetResponse>}
 */
const methodDescriptor_ZetabaseProvider_GetData = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/GetData',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.TableGet,
  proto.zbprotocol.TableGetResponse,
  /**
   * @param {!proto.zbprotocol.TableGet} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.TableGetResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.TableGet,
 *   !proto.zbprotocol.TableGetResponse>}
 */
const methodInfo_ZetabaseProvider_GetData = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.TableGetResponse,
  /**
   * @param {!proto.zbprotocol.TableGet} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.TableGetResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.TableGet} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.TableGetResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.TableGetResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.getData =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/GetData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_GetData,
      callback);
};


/**
 * @param {!proto.zbprotocol.TableGet} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.TableGetResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.getData =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/GetData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_GetData);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.TableQuery,
 *   !proto.zbprotocol.TableGetResponse>}
 */
const methodDescriptor_ZetabaseProvider_QueryData = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/QueryData',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.TableQuery,
  proto.zbprotocol.TableGetResponse,
  /**
   * @param {!proto.zbprotocol.TableQuery} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.TableGetResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.TableQuery,
 *   !proto.zbprotocol.TableGetResponse>}
 */
const methodInfo_ZetabaseProvider_QueryData = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.TableGetResponse,
  /**
   * @param {!proto.zbprotocol.TableQuery} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.TableGetResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.TableQuery} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.TableGetResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.TableGetResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.queryData =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/QueryData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_QueryData,
      callback);
};


/**
 * @param {!proto.zbprotocol.TableQuery} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.TableGetResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.queryData =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/QueryData',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_QueryData);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.NewSubIdentityRequest,
 *   !proto.zbprotocol.NewIdentityResponse>}
 */
const methodDescriptor_ZetabaseProvider_CreateUser = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/CreateUser',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.NewSubIdentityRequest,
  proto.zbprotocol.NewIdentityResponse,
  /**
   * @param {!proto.zbprotocol.NewSubIdentityRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.NewIdentityResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.NewSubIdentityRequest,
 *   !proto.zbprotocol.NewIdentityResponse>}
 */
const methodInfo_ZetabaseProvider_CreateUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.NewIdentityResponse,
  /**
   * @param {!proto.zbprotocol.NewSubIdentityRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.NewIdentityResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.NewSubIdentityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.NewIdentityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.NewIdentityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.createUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/CreateUser',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_CreateUser,
      callback);
};


/**
 * @param {!proto.zbprotocol.NewSubIdentityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.NewIdentityResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.createUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/CreateUser',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_CreateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.PermissionsEntry,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_SetPermission = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/SetPermission',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.PermissionsEntry,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.PermissionsEntry} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.PermissionsEntry,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_SetPermission = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.PermissionsEntry} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.PermissionsEntry} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.setPermission =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/SetPermission',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_SetPermission,
      callback);
};


/**
 * @param {!proto.zbprotocol.PermissionsEntry} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.setPermission =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/SetPermission',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_SetPermission);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.AuthenticateUser,
 *   !proto.zbprotocol.AuthenticateUserResponse>}
 */
const methodDescriptor_ZetabaseProvider_LoginUser = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/LoginUser',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.AuthenticateUser,
  proto.zbprotocol.AuthenticateUserResponse,
  /**
   * @param {!proto.zbprotocol.AuthenticateUser} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.AuthenticateUserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.AuthenticateUser,
 *   !proto.zbprotocol.AuthenticateUserResponse>}
 */
const methodInfo_ZetabaseProvider_LoginUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.AuthenticateUserResponse,
  /**
   * @param {!proto.zbprotocol.AuthenticateUser} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.AuthenticateUserResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.AuthenticateUser} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.AuthenticateUserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.AuthenticateUserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.loginUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/LoginUser',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_LoginUser,
      callback);
};


/**
 * @param {!proto.zbprotocol.AuthenticateUser} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.AuthenticateUserResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.loginUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/LoginUser',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_LoginUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.ListTablesRequest,
 *   !proto.zbprotocol.ListTablesResponse>}
 */
const methodDescriptor_ZetabaseProvider_ListTables = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/ListTables',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.ListTablesRequest,
  proto.zbprotocol.ListTablesResponse,
  /**
   * @param {!proto.zbprotocol.ListTablesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ListTablesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.ListTablesRequest,
 *   !proto.zbprotocol.ListTablesResponse>}
 */
const methodInfo_ZetabaseProvider_ListTables = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ListTablesResponse,
  /**
   * @param {!proto.zbprotocol.ListTablesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ListTablesResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.ListTablesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ListTablesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ListTablesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.listTables =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListTables',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListTables,
      callback);
};


/**
 * @param {!proto.zbprotocol.ListTablesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ListTablesResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.listTables =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListTables',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListTables);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.ListKeysRequest,
 *   !proto.zbprotocol.ListKeysResponse>}
 */
const methodDescriptor_ZetabaseProvider_ListKeys = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/ListKeys',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.ListKeysRequest,
  proto.zbprotocol.ListKeysResponse,
  /**
   * @param {!proto.zbprotocol.ListKeysRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ListKeysResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.ListKeysRequest,
 *   !proto.zbprotocol.ListKeysResponse>}
 */
const methodInfo_ZetabaseProvider_ListKeys = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ListKeysResponse,
  /**
   * @param {!proto.zbprotocol.ListKeysRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ListKeysResponse.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.ListKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ListKeysResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ListKeysResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.listKeys =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListKeys',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListKeys,
      callback);
};


/**
 * @param {!proto.zbprotocol.ListKeysRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ListKeysResponse>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.listKeys =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/ListKeys',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_ListKeys);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zbprotocol.DeleteSystemObjectRequest,
 *   !proto.zbprotocol.ZbError>}
 */
const methodDescriptor_ZetabaseProvider_DeleteObject = new grpc.web.MethodDescriptor(
  '/zbprotocol.ZetabaseProvider/DeleteObject',
  grpc.web.MethodType.UNARY,
  proto.zbprotocol.DeleteSystemObjectRequest,
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.DeleteSystemObjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zbprotocol.DeleteSystemObjectRequest,
 *   !proto.zbprotocol.ZbError>}
 */
const methodInfo_ZetabaseProvider_DeleteObject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zbprotocol.ZbError,
  /**
   * @param {!proto.zbprotocol.DeleteSystemObjectRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zbprotocol.ZbError.deserializeBinary
);


/**
 * @param {!proto.zbprotocol.DeleteSystemObjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zbprotocol.ZbError)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zbprotocol.ZbError>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zbprotocol.ZetabaseProviderClient.prototype.deleteObject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/DeleteObject',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_DeleteObject,
      callback);
};


/**
 * @param {!proto.zbprotocol.DeleteSystemObjectRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zbprotocol.ZbError>}
 *     A native promise that resolves to the response
 */
proto.zbprotocol.ZetabaseProviderPromiseClient.prototype.deleteObject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zbprotocol.ZetabaseProvider/DeleteObject',
      request,
      metadata || {},
      methodDescriptor_ZetabaseProvider_DeleteObject);
};


module.exports = proto.zbprotocol;

