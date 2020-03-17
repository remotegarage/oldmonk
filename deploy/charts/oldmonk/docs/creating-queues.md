This Helm Chart just create a Operator resource, this example model is used to create a QueueAutoScaler resource to every `queue` declared on values, this template/values can be added yo a existant helm chart.

### Create a new template file:
```yaml
{{- if .Values.queues }}	
{{- range $queue := .Values.queues }}	
---
apiVersion: v1  
kind: Secret  
metadata: 
  name: {{ template "oldmonk.fullname" $ }}-queue-{{ $queue.name | replace "_" "-"}}  
  namespace: {{ $.Release.Namespace }}  
type: Opaque  
data: 
  URI: {{ $queue.uri | b64enc | quote }}  
---
apiVersion: oldmonk.evalsocket.in/v1	
kind: QueueAutoScaler	
metadata:	
  name: {{ template "oldmonk.fullname" $ }}-queue-{{ $queue.name | replace "_" "-"}}	
  namespace: {{ $.Release.Namespace }}	
spec:	
  type: {{ $queue.type | quote }}	
  policy: {{ $queue.policy | quote }}	
  {{- if $queue.targetMessagesPerWorker }}	
  targetMessagesPerWorker: {{ $queue.targetMessagesPerWorker }}	
  {{- end }}	
  {{- if $queue.option }}	
  option:	
    {{- $queue.option | toYaml | trimSuffix "\n" | nindent 4 }}	
  {{- end }}	
  minPods: {{ $queue.minPods }}	
  maxPods: {{ $queue.maxPods }}	
  scaleDown:	
    amount: {{ $queue.scaleDown.amount }}	
    threshold: {{ $queue.scaleDown.threshold }}	
  scaleUp:	
    amount: {{ $queue.scaleUp.amount }}	
    threshold: {{ $queue.scaleUp.threshold }}	
  deployment: {{ $queue.deployment }}	
  {{- if $queue.labels }}	
  labels:	
    {{- $queue.labels | toYaml | trimSuffix "\n" | nindent 4 }}	
  {{- end }}	
  autopilot: {{ $queue.autopilot }}	
  secrets: "{{ template "oldmonk.fullname" $ }}-queue-{{ $queue.name | replace "_" "-"}}"	
{{ end }}	
{{ end }}
```

### Add to values:
```yaml
queues:
- name: test 
  uri: "https://sqs.us-east-1.amazonaws.com/111111111111/test" 
  type: "SQS"  
  policy: "THRESOLD" 
  option:  
    region: 'us-east-1'  
  minPods: 0 
  maxPods: 20  
  scaleDown: 
    amount: 5  
    threshold: 10  
  scaleUp: 
    amount: 5  
    threshold: 10  
  deployment: 'worker' 
  labels:  
    app: worker  
  autopilot: false
```

### What will happen?
The `template` has a loop that will create a QueueAutoScaler resource to every `queue` declared on values.yaml.