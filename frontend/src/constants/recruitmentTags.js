export const recruitmentTagGroups = [
  {
    key: 'profession',
    label: '职业',
    tags: ['近卫', '狙击', '重装', '医疗', '辅助', '术师', '特种', '先锋'],
  },
  {
    key: 'position',
    label: '部署位置',
    tags: ['近战位', '远程位'],
  },
  {
    key: 'traits',
    label: '特性标签',
    tags: ['控场', '爆发', '治疗', '支援', '费用回复', '输出', '生存', '群攻', '防护', '减速', '削弱', '快速复活', '位移', '召唤', '支援机械', '元素'],
  },
  {
    key: 'seniority',
    label: '干员资历',
    tags: ['新手', '资深干员', '高级资深干员'],
  },
]

export const recruitmentTagSet = new Set(
  recruitmentTagGroups.flatMap((group) => group.tags),
)
