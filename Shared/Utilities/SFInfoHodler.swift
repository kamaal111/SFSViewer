//
//  SFInfo.swift
//  SFSViewer
//
//  Created by Kamaal M Farah on 23/12/2021.
//

import Foundation

public struct SFInfoHodler {
    public let items: [SFInfo]

    public init() throws {
        let supportedVersionsDict: [String: SFVersions] = try Self.decodeFile(withName: "supported_versions")
        let names: [SFName] = try Self.decodeFile(withName: "names")

        self.items = names.compactMap({
            guard let supportedVersions = supportedVersionsDict[$0.releaseYear] else { return nil }
            return SFInfo(name: $0.name, supportedVersions: supportedVersions)
        })
    }

    public enum InitializerErrors: Error {
        case resourceNotFound(name: String)
        case encodingError(error: Error)
        case decodingError(error: Error)
    }

    private static func decodeFile<T: Decodable>(withName name: String) throws -> T {
        guard let url = Bundle.main.url(forResource: name, withExtension: "json")
        else { throw InitializerErrors.resourceNotFound(name: name) }

        let data: Data
        do {
            data = try Data(contentsOf: url)
        } catch {
            throw InitializerErrors.encodingError(error: error)
        }

        let decodedData: T
        do {
            decodedData = try JSONDecoder().decode(T.self, from: data)
        } catch {
            throw InitializerErrors.decodingError(error: error)
        }

        return decodedData
    }
}

extension SFInfoHodler: CustomStringConvertible {
    public var description: String {
        "SFInfoHodler(items: \(items.count) items)"
    }
}

public struct SFInfo: Hashable {
    public let name: String
    public let supportedVersions: SFVersions

    public init(name: String, supportedVersions: SFVersions) {
        self.name = name
        self.supportedVersions = supportedVersions
    }
}

struct SFName: Codable, Hashable {
    let name: String
    let releaseYear: String

    init(name: String, releaseYear: String) {
        self.name = name
        self.releaseYear = releaseYear
    }

    enum CodingKeys: String, CodingKey {
        case name
        case releaseYear = "release_year"
    }
}

public struct SFVersions: Codable, Hashable {
    public let iOS: String
    public let macOS: String
    public let tvOS: String
    public let watchOS: String

    public init(iOS: String, macOS: String, tvOS: String, watchOS: String) {
        self.iOS = iOS
        self.macOS = macOS
        self.tvOS = tvOS
        self.watchOS = watchOS
    }
}
