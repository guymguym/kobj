KOBJ_URL=https://localhost:8443/apis/kobj.io/v1alpha1/namespaces/default/kobjs
curl -k $KOBJ_URL
curl -k $KOBJ_URL -X POST -H "content-type: application/json" -d '{"metadata":{"name":"ggg","namespace":"default"},"value":"gagagagaga"}'
