package kafka

import "github.com/Hvaekar/med-account/pkg/broker"

var topics = map[string]string{
	broker.AccountCreateKey: "account.create",
	broker.AccountDeleteKey: "account.delete",
	broker.AccountGetKey:    "account.get",
	broker.AccountLoginKey:  "account.enter",
	broker.AccountLogoutKey: "account.exit",
	broker.AccountUpdateKey: "account.update",

	broker.FileAddKey:    "account_file.add",
	broker.FileDeleteKey: "account_file.delete",
	broker.FileGetKey:    "account_file.get",
	broker.FileUpdateKey: "account_file.update",

	broker.AddressAddKey:    "account_address.add",
	broker.AddressDeleteKey: "account_address.delete",
	broker.AddressGetKey:    "account_address.get",
	broker.AddressUpdateKey: "account_address.update",

	broker.EmailAddKey:      "account_email.add",
	broker.EmailDeleteKey:   "account_email.delete",
	broker.EmailGetKey:      "account_email.get",
	broker.EmailUpdateKey:   "account_email.update",
	broker.EmailVerifiedKey: "account_email.verified",
	broker.EmailVerifyKey:   "account_email.verify",

	broker.PhoneAddKey:      "account_phone.add",
	broker.PhoneDeleteKey:   "account_phone.delete",
	broker.PhoneGetKey:      "account_phone.get",
	broker.PhoneUpdateKey:   "account_phone.update",
	broker.PhoneVerifiedKey: "account_phone.verified",
	broker.PhoneVerifyKey:   "account_phone.verify",

	broker.LanguageAddKey:    "account_language.add",
	broker.LanguageDeleteKey: "account_language.delete",
	broker.LanguageGetKey:    "account_language.get",
	broker.LanguageUpdateKey: "account_language.update",

	broker.ProfileGetKey: "account_profile.get",

	broker.PatientDeleteKey: "patient.delete",
	broker.PatientGetKey:    "patient.get",
	broker.PatientUpdateKey: "patient.update",

	broker.PatientVerifiedKey: "patient.verified",
	broker.PatientSelectedKey: "patient.selected",

	broker.PatientAdminAddKey:    "patient_admin.add",
	broker.PatientAdminDeleteKey: "patient_admin.delete",
	broker.PatientAdminGetKey:    "patient_admin.get",
	broker.PatientAdminUpdateKey: "patient_admin.update",

	broker.MetalComponentAddKey:    "patient_metal_component.add",
	broker.MetalComponentDeleteKey: "patient_metal_component.delete",
	broker.MetalComponentGetKey:    "patient_metal_component.get",
	broker.MetalComponentUpdateKey: "patient_metal_component.update",

	broker.SpecialistAddKey:    "specialist.add",
	broker.SpecialistGetKey:    "specialist.get",
	broker.SpecialistUpdateKey: "specialist.update",

	broker.SpecializationAddKey:    "specialist_specialization.add",
	broker.SpecializationDeleteKey: "specialist_specialization.delete",
	broker.SpecializationGetKey:    "specialist_specialization.get",
	broker.SpecializationUpdateKey: "specialist_specialization.update",

	broker.AssociationAddKey:    "specialist_association.add",
	broker.AssociationDeleteKey: "specialist_association.delete",
	broker.AssociationGetKey:    "specialist_association.get",
	broker.AssociationUpdateKey: "specialist_association.update",

	broker.EducationAddKey:    "specialist_education.add",
	broker.EducationDeleteKey: "specialist_education.delete",
	broker.EducationGetKey:    "specialist_education.get",
	broker.EducationUpdateKey: "specialist_education.update",

	broker.ExperienceAddKey:    "specialist_experience.add",
	broker.ExperienceDeleteKey: "specialist_experience.delete",
	broker.ExperienceGetKey:    "specialist_experience.get",
	broker.ExperienceUpdateKey: "specialist_experience.update",

	broker.PatentAddKey:    "specialist_patent.add",
	broker.PatentDeleteKey: "specialist_patent.delete",
	broker.PatentGetKey:    "specialist_patent.get",
	broker.PatentUpdateKey: "specialist_patent.update",

	broker.PublicationLinkAddKey:    "specialist_publication_link.add",
	broker.PublicationLinkDeleteKey: "specialist_publication_link.delete",
	broker.PublicationLinkGetKey:    "specialist_publication_link.get",
	broker.PublicationLinkUpdateKey: "specialist_publication_link.update",
}

func (b *MessageBroker) GetTopic(name string) string {
	topic, ok := topics[name]
	if !ok {
		return ""
	}

	return topic
}
