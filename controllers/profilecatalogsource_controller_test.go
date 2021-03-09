/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
)

const (
	catalogSourceName = "test-catalog"
)

var _ = Describe("ProfileCatalogSourceReconciler", func() {
	var (
		catalogSource *profilesv1alpha1.ProfileCatalogSource
	)

	Context("when a ProfileCatalogSource is created", func() {
		BeforeEach(func() {
			ctx := context.Background()
			catalogSource = &profilesv1alpha1.ProfileCatalogSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ProfileCatalogSource",
					APIVersion: "profiles.weave.works/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      catalogSourceName,
					Namespace: testNamespace,
				},
				Spec: profilesv1alpha1.ProfileCatalogSourceSpec{
					EmbeddedProfiles: []profilesv1alpha1.Profile{
						{
							Name:        "test-profile",
							URL:         "https://github.com/testing/testing.git",
							Version:     "v0.0.1",
							Description: "This is a test profile",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, catalogSource)).To(Succeed())
		})

		AfterEach(func() {
			testProfiles.reset()
			Expect(k8sClient.Delete(context.TODO(), catalogSource)).To(Succeed())
		})

		It("records the discovered profiles", func() {
			Eventually(func() []profilesv1alpha1.Profile {
				return testProfiles.Saved
			}, timeout, time.Millisecond*500).Should(Equal(catalogSource.Spec.EmbeddedProfiles))
		})
	})
})
