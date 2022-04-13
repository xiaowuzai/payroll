FROM alpine:3.12

WORKDIR $GOPATH/src/app
ADD payroll $GOPATH/src/app/payroll/
ADD payroll/configs/ $GOPATH/src/app/configs/

CMD payroll/payroll