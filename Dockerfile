#
#   Author: Rohith (gambol99@gmail.com)
#   Date: 2015-02-19 16:39:24 +0000 (Thu, 19 Feb 2015)
#
#  vim:ts=2:sw=2:et
#
FROM progrium/busybox
MAINTAINER Rohith <gambol99@gmail.com>

ADD stage/fabric /bin/fabric
RUN opkg-install bash
RUN chmod +x /bin/fabric

EXPOSE 1022 7369
ENTRYPOINT [ "/bin/fabric" ]
