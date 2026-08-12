# Changelog

## [1.2.3](https://github.com/na1tto/go-social/compare/v1.2.2...v1.2.3) (2026-08-12)


### Bug Fixes

* fix typo in Level field in roles.go ([6a300b0](https://github.com/na1tto/go-social/commit/6a300b0d2bb57619343d39adbd96df5aa238e099))

## [1.2.2](https://github.com/na1tto/go-social/compare/v1.2.1...v1.2.2) (2026-08-11)


### Bug Fixes

* fixed missing certificate extension ([9b94b0d](https://github.com/na1tto/go-social/commit/9b94b0d7c71f78275d753ff25af3f02b80ab5ab5))

## [1.2.1](https://github.com/na1tto/go-social/compare/v1.2.0...v1.2.1) (2026-08-10)


### Bug Fixes

* fixed missing quotation marks in dockerfile cmd path ([b8d2b5a](https://github.com/na1tto/go-social/commit/b8d2b5a63ce982a1217d23e301eb3a9572f2c19e))

## [1.2.0](https://github.com/na1tto/go-social/compare/v1.1.0...v1.2.0) (2026-08-09)


### Features

* added dockerfile ([600ac48](https://github.com/na1tto/go-social/commit/600ac482494ed86f9f4ce29db67ee9bd9850a045))

## [1.1.0](https://github.com/na1tto/go-social/compare/v1.0.0...v1.1.0) (2026-08-09)


### Features

* added api version automation ([1facb5d](https://github.com/na1tto/go-social/commit/1facb5d053da8701c0576fd3d0ac2e005f0cb7da))

## 1.0.0 (2026-08-09)


### Features

* add initial README with project overview and setup instructions ([157068b](https://github.com/na1tto/go-social/commit/157068bb94a629b8b7343397e6f334af0bd43b27))
* added activation invitation table and user activated status column to the db ([c2817fb](https://github.com/na1tto/go-social/commit/c2817fb7538d75bb3439d5e3d43e8e7f74a5caa2))
* added automation workflow ([c51a1f1](https://github.com/na1tto/go-social/commit/c51a1f15e68a6a47a8e8520fe9dd57c1b718afeb))
* added basic auth to the /health endpoint ([ac3371a](https://github.com/na1tto/go-social/commit/ac3371ac69a7f335f7f97047a8b3bf765fd7cf43))
* added comments creation api endpoint ([fdc55d8](https://github.com/na1tto/go-social/commit/fdc55d8d48a9b02fc861ad6b7c6d0e9749ab3d32))
* added comments for posts ([fdea9f2](https://github.com/na1tto/go-social/commit/fdea9f27def3b7c0d1809fcf028e90113c1d08f8))
* added cors ([7c7ac40](https://github.com/na1tto/go-social/commit/7c7ac40b9472d1627597e15c564d5ef148025a8f))
* added filtering in the feed endpoint ([c18247f](https://github.com/na1tto/go-social/commit/c18247faaa19398a0cbef00f9994aad3d04257e6))
* added followers logic to the application, implemented middleware for fetching users data ([41fe318](https://github.com/na1tto/go-social/commit/41fe3186c7823ddfe7570b87dc951a6afe02a457))
* added graceful shutdown ([109e759](https://github.com/na1tto/go-social/commit/109e75998c5ecaa1cb5685ebc7f672fc13bde5f3))
* added inital auth logic for user creation and invitation via email ([ae8d874](https://github.com/na1tto/go-social/commit/ae8d874d4891f5459e200c4bea9ad9c790798ca5))
* added initial stateless auth logic ([fec1077](https://github.com/na1tto/go-social/commit/fec1077a580bc851cd457e2d9d2c1adf169482b9))
* added mailtrap email support ([5900d39](https://github.com/na1tto/go-social/commit/5900d39e3c5d28b8a6586cede2ce69c783b942f4))
* added metrics endpoint and removed cors allowed origin for now ([934b0c2](https://github.com/na1tto/go-social/commit/934b0c213a252cf8d8d39933eddb1e14785319a8))
* added mocking tests ([7b703eb](https://github.com/na1tto/go-social/commit/7b703eb60a366bd3f39cbeb458bf402658b9b7f4))
* added password validation in the authentication/token endpoint ([80dd57e](https://github.com/na1tto/go-social/commit/80dd57ec2242a4ff86295df7c33986d2b0dc3cfe))
* added rate limiter fixed window implementation ([f98afdd](https://github.com/na1tto/go-social/commit/f98afddfa5392a59924c55639cdd6ad431464dc5))
* added redis cache to endpoints ([7d3fe79](https://github.com/na1tto/go-social/commit/7d3fe79822c961339b38e4e6aa9192e99e38fc30))
* added redis client abstraction ([5a6e280](https://github.com/na1tto/go-social/commit/5a6e28085b54898e7b9f611677ffc4aff14df55a))
* added roles based authorization ([0c0a72d](https://github.com/na1tto/go-social/commit/0c0a72d3db5b72842c2771ab5c5cc6dc617ad5f6))
* added sendmail for account activation via email ([d5d3e73](https://github.com/na1tto/go-social/commit/d5d3e736e8cb036a46af532e10407a1c267bbf2d))
* added structure loggin using uber's zap package ([06e1642](https://github.com/na1tto/go-social/commit/06e1642f134a0c32daa5ac792b63da0a8c1a031c))
* added swagger docs ([0ba2f9a](https://github.com/na1tto/go-social/commit/0ba2f9a7de65adcc8482b47e306ff73da13909b9))
* added the authentication endopoint to be worked on ([cf260df](https://github.com/na1tto/go-social/commit/cf260df08bd1d827428aa7d5a627c5884ecb2b2b))
* added the frontend initial files for future uses ([1d16d92](https://github.com/na1tto/go-social/commit/1d16d929f26edf6cc79ec6ca9b323a31d99fbcce))
* added token validation and fixed some typos ([1a56334](https://github.com/na1tto/go-social/commit/1a56334c329ea729292944a044aa21ab60e01801))
* added user confirmation page ([d0dfa66](https://github.com/na1tto/go-social/commit/d0dfa66b44a0cd66229f3843fd656dd0fb5071e7))
* **database:** added roles table and updated all users to user role ([8c66397](https://github.com/na1tto/go-social/commit/8c663976a77e18341df0771bfa0e41a5ee669f7b))
* Docker postgres container implemented, migrations implemented, makefile created for migrations ([88d33b8](https://github.com/na1tto/go-social/commit/88d33b8d77606952247e6522faf3e5fe95156042))
* feed endpoint and query added ([56a1593](https://github.com/na1tto/go-social/commit/56a15931c43eb9c444a6050e070debefcbe551a2))
* finished user creation and invitation endpoints ([3bea299](https://github.com/na1tto/go-social/commit/3bea29926baba4776b16df8128fb301b6b52ca3b))
* get user endpoint implemented ([35f8886](https://github.com/na1tto/go-social/commit/35f88861d7e25b618dcb1f379065fc5585dbab69))
* internal error handling in the api layer implemented ([c904d83](https://github.com/na1tto/go-social/commit/c904d839c8b329ee41734cb3d48401d900925c16))
* json response and writers for the CRUD operations ([f2db4e6](https://github.com/na1tto/go-social/commit/f2db4e688f7230a4a2ab004029d396b28e623fff))
* middleware for posts fetching added, delete and update functions implemented, comments query fixed ([3af6ff0](https://github.com/na1tto/go-social/commit/3af6ff04911dfc4a777c940153d5c89aedcb5b60))
* more registerUserHandler logic, seed not working in this commit ([ddd6c4f](https://github.com/na1tto/go-social/commit/ddd6c4f7e537922968fde048ff58bf20b28e017a))
* pagination and simple sorting added ([85818cc](https://github.com/na1tto/go-social/commit/85818cca1ebeb1f3886023c4fe501ceec5abca28))
* post create operation ([7b0d535](https://github.com/na1tto/go-social/commit/7b0d535750566d7c5e70c72b8024a166b7c49688))
* post GetById method added ([d097e38](https://github.com/na1tto/go-social/commit/d097e38af5b7d5257980b0b68b52698092d3071f))
* release please script ([7ad924c](https://github.com/na1tto/go-social/commit/7ad924c3c22ee38afc0293de785ea73bc76c46d4))
* Repository pattern mocked ([c89f1d3](https://github.com/na1tto/go-social/commit/c89f1d3b015d998c6df9836e0de59c26c4b66d67))
* User and Post create methods done ([0c891ea](https://github.com/na1tto/go-social/commit/0c891ea0c19d06d9fc4446f8c0adb829327cda26))
* version column in the posts table to handle concurrent update requests in the data layer ([2a55c2f](https://github.com/na1tto/go-social/commit/2a55c2fe8e5063608729d26573653666a8da2a16))


### Bug Fixes

* added return on comment handler 400 response ([a7b53dc](https://github.com/na1tto/go-social/commit/a7b53dca5860f09275c8041a13c7599b4b62252f))
* auth endpoint adjusted to outside users group, password is hashed in the create user fn and seed tx added ([4d3555a](https://github.com/na1tto/go-social/commit/4d3555a8df8c581ada5fead15eaa815d90fefc7e))
* changed created_at field to global) ([ccf917b](https://github.com/na1tto/go-social/commit/ccf917b8974880c998e56796224956b463342758))
* default password value setter added for the db seeding ([1a4d250](https://github.com/na1tto/go-social/commit/1a4d25087ea65e15646d6e1f40b267c5a487832a))
* fixed error response in getUser method ([703c92a](https://github.com/na1tto/go-social/commit/703c92a71bd94b35641fd9f377210a5540b15313))
* fixed follow and unfollow parameter parsing error, removed post id mocking, standarized error response formats ([20689c4](https://github.com/na1tto/go-social/commit/20689c4175e02e9083fca612ba1ff151304ff4ba))
* fixed GetById query ([8a7cb01](https://github.com/na1tto/go-social/commit/8a7cb014dabd14fcd28f46f7bf8190bce1344801))
* fixed some todos in the codebase ([40a4319](https://github.com/na1tto/go-social/commit/40a431930ddb5957be4a3a02652525de7942d315))
* fixed typo in contacts ([e796b0e](https://github.com/na1tto/go-social/commit/e796b0e29080f726402bc6cdfd4243185cb6ea93))
* fixed undefined type in posts.go ([b13dc8f](https://github.com/na1tto/go-social/commit/b13dc8fdb11eab97b2e6d2f443a4a3b60d33e402))
* removed unused middleware, unused err and sotre conflict status declaration ([d561875](https://github.com/na1tto/go-social/commit/d561875516eabad1b5c86a31e0e042660c52aba4))
* removed user id mock from the feed controller ([9a30ec9](https://github.com/na1tto/go-social/commit/9a30ec9ec2e61cfca7d8d4e99a1109ff3fa5a41d))
* renamed gen-docs duplicate ([2a22dac](https://github.com/na1tto/go-social/commit/2a22dac35d42e46a3417661c507c499a2230598d))
* user creation fixed with default role ([148026f](https://github.com/na1tto/go-social/commit/148026f4c13b01ce45d624e5a775ec5385b2ebed))


### Performance Improvements

* added indexes for full-text search ([22198c9](https://github.com/na1tto/go-social/commit/22198c94359391d9bb4da07320886109203d953f))
