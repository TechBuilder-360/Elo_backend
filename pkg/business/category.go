package business

var category []string

func init() {
	category = []string{
		"Agriculture",
		"Automotive",
		"Beauty & Personal Care",
		"Construction",
		"Consulting",
		"Education",
		"Energy & Utilities",
		"Entertainment",
		"Finance",
		"Food & Beverage",
		"Government & Public Sector",
		"Healthcare",
		"Hospitality",
		"Information Technology",
		"Logistics & Transportation",
		"Manufacturing",
		"Marketing & Advertising",
		"Media & Publishing",
		"Mining & Natural Resources",
		"Nonprofit & NGO",
		"Professional Services",
		"Real Estate",
		"Retail & E-commerce",
		"Security Services",
		"Telecommunications",
		"Tourism & Travel",
		"Wholesale & Distribution",
		"Religious Organization",
		"Other",
	}
}

func GetCategory() []string {
	return category
}
